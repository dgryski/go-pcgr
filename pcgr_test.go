package pcgr

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/lazybeaver/xorshift"
)

/*
#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>

typedef struct {
    uint64_t state;
    uint64_t inc;
} pcg32_random_t;

uint32_t pcg32_random_r(pcg32_random_t * rng)
{
    uint64_t oldstate = rng->state;
    rng->state = oldstate * 6364136223846793005ULL + (rng->inc | 1);
    uint32_t xorshifted = ((oldstate >> 18u) ^ oldstate) >> 27u;
    uint32_t rot = oldstate >> 59u;
    return (xorshifted >> rot) | (xorshifted << ((-rot) & 31));
}

int main(int argc, char *argv[])
{

    pcg32_random_t rng = { 0x0ddc0ffeebadf00dULL, 0xcafebabe };

    int i;

    for (i = 0; i < 10000; i++) {
        printf("%u\n", pcg32_random_r(&rng));
    }

    return 0;
}
*/

func TestGenerate(t *testing.T) {

	rnd := Rand{0x0ddc0ffeebadf00d, 0xcafebabe}

	// generated from the above reference C code
	f, err := os.Open("testdata/numbers.txt")
	if err != nil {
		t.Fatalf("unable to open data set: %v ", err)
	}

	scanner := bufio.NewScanner(f)

	var line int
	for scanner.Scan() {
		n := rnd.Next()
		want, err := strconv.Atoi(scanner.Text())
		if err != nil {
			t.Fatalf("unable to parse data line %d: %v\n", line, err)
		}
		line++
		if n != uint32(want) {
			t.Fatalf("rng mismatch round %d: got %d want %d\n", line, n, uint32(want))
		}
	}
}

func TestAdvance(t *testing.T) {

	rnd := Rand{0x0ddc0ffeebadf00d, 0xcafebabe}

	var ints []uint32

	for i := 0; i < 10; i++ {
		ints = append(ints, rnd.Next())
	}

	rnd.Advance(-10)

	for i := 0; i < 10; i++ {
		if n := rnd.Next(); n != ints[i] {
			t.Errorf("advance failed: step %d = %d, want %d\n", i, n, ints[i])
		}
	}

	tmp := rnd

	for i := 0; i < 100; i++ {
		tmp.Next()
	}

	rnd.Advance(100)

	for i := 0; i < 10; i++ {
		if got, want := rnd.Next(), tmp.Next(); got != want {
			t.Errorf("advance failed: step %d = %d, want %d\n", i, got, want)
		}
	}
}

func TestCompat(t *testing.T) {

	// from pcg32-demo
	output := []struct {
		numbers string
		coins   string
		rolls   string
		cards   string
	}{
		{
			"0xa15c02b7 0x7b47f409 0xba1d3330 0x83d2f293 0xbfa4784b 0xcbed606e",
			"HHTTTHTHHHTHTTTHHHHHTTTHHHTHTHTHTTHTTTHHHHHHTTTTHHTTTTTHTTTTTTTHT",
			"3 4 1 1 2 2 3 2 4 3 2 4 3 3 5 2 3 1 3 1 5 1 4 1 5 6 4 6 6 2 6 3 3",
			"Qd Ks 6d 3s 3d 4c 3h Td Kc 5c Jh Kd Jd As 4s 4h Ad Th Ac Jc 7s Qs 2s 7h Kh 2d 6c Ah 4d Qh 9h 6s 5s 2c 9c Ts 8d 9s 3c 8c Js 5d 2h 6h 7d 8s 9d 5h 8h Qc 7c Tc",
		},
		{
			"0x74ab93ad 0x1c1da000 0x494ff896 0x34462f2f 0xd308a3e5 0x0fa83bab",
			"HHHHHHHHHHTHHHTHTHTHTHTTTTHHTTTHHTHHTHTTHHTTTHHHHHHTHTTHTHTTTTTTT",
			"5 1 1 3 3 2 4 5 3 2 2 6 4 3 2 4 2 4 3 2 3 6 3 2 3 4 2 4 1 1 5 4 4",
			"7d 2s 7h Td 8s 3c 3d Js 2d Tc 4h Qs 5c 9c Th 2c Jc Qd 9d Qc 7s 3s 5s 6h 4d Jh 4c Ac 4s 5h 5d Kc 8h 8d Jd 9s Ad 6s 6c Kd 2h 3h Kh Ts Qh 9h 6d As 7c Ks Ah 8c",
		},
		{
			"0x39af5f9f 0x04196b18 0xc3c3eb28 0xc076c60c 0xc693e135 0xf8f63932",
			"HTTHHTTTTTHTTHHHTHTTHHTTHTHHTHTHTTTTHHTTTHHTHHTTHTTHHHTHHHTHTTTHT",
			"5 1 5 3 2 2 4 5 3 3 1 3 4 6 3 2 3 4 2 2 3 1 5 2 4 6 6 4 2 4 3 3 6",
			"Kd Jh Kc Qh 4d Qc 4h 9d 3c Kh Qs 8h 5c Jd 7d 8d 3h 7c 8s 3s 2h Ks 9c 9h 2c 8c Ad 7s 4s 2s 5h 6s 4c Ah 7h 5s Ac 3d 5d Qd As Tc 6h 9s 2d 6c 6d Td Jc Ts Th Js",
		},
		{
			"0x55ce6851 0x97a7726d 0x17e10815 0x58007d43 0x962fb148 0xb9bb55bd",
			"HHTHHTTTTHTHHHHHTTHHHTTTHHTHTHTHTHHTTHTHHHHHHTHHTHHTHHTTTTHHTHHTT",
			"6 6 3 2 3 4 2 6 4 2 6 3 2 3 5 5 3 4 4 6 6 2 6 5 4 4 6 1 6 1 3 6 5",
			"Qd 8h 5d 8s 8d Ts 7h Th Qs Js 7s Kc 6h 5s 4d Ac Jd 7d 7c Td 2c 6s 5h 6d 3s Kd 9s Jh Kh As Ah 9h 3c Qh 9c 2d Tc 9d 2s 3d Ks 4h Qc Ad Jc 8c 2h 3h 4s 4c 5c 6c",
		},
		{
			"0xfcef7cd6 0x1b488b5a 0xd0daf7ea 0x1d9a70f7 0x241a37cf 0x9a3857b7",
			"HHHHTHHTTHTTHHHTTTHHTHTHTTTTHTTHTHTTTHHHTHTHTTHTTHTHHTHTHHHTHTHTT",
			"5 4 1 2 6 1 3 1 5 6 3 6 2 1 4 4 5 2 1 5 6 5 6 4 4 4 5 2 6 4 3 5 6",
			"4d 9s Qc 9h As Qs 7s 4c Kd 6h 6s 2c 8c 5d 7h 5h Jc 3s 7c Jh Js Ks Tc Jd Kc Th 3h Ts Qh Ad Td 3c Ah 2d 3d 5c Ac 8s 5s 9c 2h 6c 6d Kh Qd 8d 7d 2s 8h 4h 9d 4s",
		},
	}

	var rnd Rand
	rnd.SeedWithState(42, 54)

	for i, tt := range output {
		nn := strings.Fields(tt.numbers)
		for j, n := range nn {
			want, _ := strconv.ParseUint(n, 0, 32)
			if got := rnd.Next(); got != uint32(want) {
				t.Errorf("failed round %d step uint32 %d: got %d want %d", i, j, got, want)
			}
		}

		for j, want := range tt.coins {
			var got rune
			if rnd.Bound(2) == 0 {
				got = 'T'
			} else {
				got = 'H'
			}

			if got != want {
				t.Errorf("failed round %d step coins %d: got %d want %d", i, j, got, want)
			}
		}

		for j, want := range strings.Fields(tt.rolls) {
			got := '0' + byte(rnd.Bound(6)) + 1

			if got != want[0] {
				t.Errorf("failed round %d step rolls %d: got %d want %d", i, j, got, want[0])
			}
		}

		cards := dealCards(&rnd)

		const (
			numSuits   = 4
			numNumbers = 13
			numCards   = 52

			cardNumber = "A23456789TJQK"
			cardSuit   = "hcds"
		)

		for j, want := range strings.Fields(tt.cards) {
			got := fmt.Sprintf("%c%c", cardNumber[cards[j]/numSuits], cardSuit[cards[j]%numSuits])
			if got != want {
				t.Errorf("failed round %d step cards %d: got %s want %s", i, j, got, want)
			}
		}
	}
}

// shuffle a deck
func dealCards(r *Rand) [52]int {

	var cards [52]int

	for i := 0; i < 52; i++ {
		cards[i] = i
	}

	for i := 52; i > 1; i-- {
		chosen := r.Bound(uint32(i))
		card := cards[chosen]
		cards[chosen] = cards[i-1]
		cards[i-1] = card
	}

	return cards
}

var total uint32 = 0

func BenchmarkPCGR(b *testing.B) {
	rnd := Rand{0x0ddc0ffeebadf00d, 0xcafebabe}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		total += rnd.Next()
	}
}

func BenchmarkPCGR64(b *testing.B) {
	rnd := Rand{0x0ddc0ffeebadf00d, 0xcafebabe}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		total += rnd.Next() + rnd.Next()
	}
}

func BenchmarkRand(b *testing.B) {

	r := rand.New(rand.NewSource(0x0ddc0ffeebadf00d))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		total += uint32(r.Int31())
	}
}
func BenchmarkXorshift(b *testing.B) {

	r := xorshift.NewXorShift64Star(0x0ddc0ffeebadf00d)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		total += uint32(r.Next())
	}
}

var _ = rand.Source(&Rand{})
