package pcgr

import (
	"bufio"
	"math/rand"
	"os"
	"strconv"
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
	output := []uint32{
		0xa15c02b7, 0x7b47f409, 0xba1d3330, 0x83d2f293, 0xbfa4784b, 0xcbed606e,
	}

	var rnd Rand
	rnd.SeedWithState(42, 54)

	for i, want := range output {
		if n := rnd.Next(); n != want {
			t.Errorf("failed step %d: got %d want %d", i, n, want)
		}
	}
}

var _ = rand.Source(&Rand{})
