package pcgr

import (
	"bufio"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
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
	}

	var rnd Rand
	rnd.SeedWithState(42, 54)

	for i, tt := range output {
		nn := strings.Fields(tt.numbers)
		for j, n := range nn {
			want, _ := strconv.ParseUint(n, 0, 32)
			if got := rnd.Next(); got != uint32(want) {
				t.Errorf("failed round %d step %d: got %d want %d", i, j, got, want)
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
	}
}

var _ = rand.Source(&Rand{})
