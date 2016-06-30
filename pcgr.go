// Package pcgr implements a small PCG random number generator
/* http://www.pcg-random.org/ */
package pcgr

// Rand is a PCG random number generator.  It implements the rand.Source interface.
type Rand struct {
	State uint64
	Inc   uint64 // must be odd
}

const defaultMultiplier64 = 6364136223846793005

// New returns an pcgr stream
func New(seed, state int64) Rand {
	var r Rand
	r.SeedWithState(seed, state)
	return r
}

// Next returns a random uint32
func (r *Rand) Next() uint32 {

	oldstate := r.State

	// Advance internal state
	r.step()

	// Calculate output function (XSH RR), uses old state for max ILP
	xorshifted := uint32(((oldstate >> 18) ^ oldstate) >> 27)
	rot := uint32(oldstate >> 59)
	return (xorshifted >> rot) | (xorshifted << ((-rot) & 31))
}

// SeedWithState sets the internal state and sequence number of the rng
func (r *Rand) SeedWithState(initstate, initseq int64) {
	r.State = 0
	r.Inc = uint64(initseq<<1) | 1 // Inc will be made odd even if initseq is not
	r.step()
	r.State += uint64(initstate)
	r.step()
}

// Seed states the internal state of the rng
func (r *Rand) Seed(seed int64) {
	r.SeedWithState(seed, 0)
}

// Int63 returns a random 63-bit integer
func (r *Rand) Int63() int64 {
	n := int64(r.Next())<<32 | int64(r.Next())
	n &= 0x7FFFFFFFFFFFFFFF
	return n
}

func (r *Rand) step() {
	r.State = r.State*defaultMultiplier64 + r.Inc
}

// Advance skips forward 'delta' steps in the stream.  Delta can be negative in which case the stream in rewound.
func (r *Rand) Advance(delta int) {

	udelta := uint64(delta)

	curMult := uint64(defaultMultiplier64)
	curPlus := r.Inc | 1

	accMult := uint64(1)
	accPlus := uint64(0)

	for udelta > 0 {
		if (udelta & 1) == 1 {
			accMult *= curMult
			accPlus = accPlus*curMult + curPlus
		}
		curPlus = (curMult + 1) * curPlus
		curMult *= curMult
		udelta /= 2
	}

	r.State = accMult*r.State + accPlus
}

// Bound returns a uniform integer 0..bound-1
func (r *Rand) Bound(bound uint32) uint32 {
	threshold := -bound % bound
	for {
		n := r.Next()
		if n >= threshold {
			return n % bound
		}
	}
}
