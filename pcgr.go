// Package pcgr implements a small PCG random number generator
/* http://www.pcg-random.org/ */
package pcgr

type Rand struct {
	State uint64
	Inc   uint64
}

const defaultMultiplier64 = 6364136223846793005

func (r *Rand) Next() uint32 {

	oldstate := r.State

	// Advance internal state
	r.step()

	// Calculate output function (XSH RR), uses old state for max ILP
	xorshifted := uint32(((oldstate >> 18) ^ oldstate) >> 27)
	rot := uint32(oldstate >> 59)
	return (xorshifted >> rot) | (xorshifted << ((-rot) & 31))
}

func (r *Rand) SeedWithState(initstate, initseq int64) {
	r.State = 0
	r.Inc = uint64(initseq<<1) | 1
	r.step()
	r.State += uint64(initstate)
	r.step()
}

func (r *Rand) Seed(seed int64) {
	r.SeedWithState(seed, 0)
}

func (r *Rand) Int63() int64 {
	n := int64(r.Next())<<32 | int64(r.Next())
	n &= 0x7FFFFFFFFFFFFFFF
	return n
}

func (r *Rand) step() {
	r.State = r.State*defaultMultiplier64 + (r.Inc | 1)
}
