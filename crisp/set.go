package crisp

import "errors"

// Set (or universe) of possible discrete `x` values
type Set struct {
	xmin float64
	xmax float64
	dx   float64
}

// NewSet builds a new Set instance
func NewSet(xmin, xmax, dx float64) (Set, error) {
	if dx == 0 {
		return Set{}, errors.New("crisp set: dx shall be > 0")
	}

	if xmin > xmax {
		return Set{}, errors.New("crisp set: xmin shall be < xmax")
	}

	return Set{
		xmin: xmin,
		xmax: xmax,
		dx:   dx,
	}, nil
}

// NewSetN builds a new Set instance with N values
func NewSetN(xmin, xmax float64, n int) (Set, error) {
	if n <= 1 {
		return Set{}, errors.New("crisp set: n shall be >= 2")
	}

	if xmin > xmax {
		return Set{}, errors.New("crisp set: xmin shall be < xmax")
	}

	return Set{
		xmin: xmin,
		xmax: xmax,
		dx:   (xmax - xmin) / float64(n-1),
	}, nil
}

// Values translates the interval into discrete increasing values
func (set Set) Values() []float64 {
	if set.dx == 0 {
		return nil
	}

	// /!\ the loop x += dx introduces a constant delta error at each step
	// Eg.: try loop from 0 to 5 with a step of 0.1
	// This error grows all over the steps

	// var result []float64
	// for x := set.xmin; x <= set.xmax; x = x + set.dx {
	// 	result = append(result, x)
	// }

	// Prefer the solution of x = min + i*dx (a delta error is still present but more acceptable)
	n := int(1 + (set.xmax-set.xmin)/set.dx)
	result := make([]float64, n)
	for i := 0; i < n; i++ {
		result[i] = set.xmin + float64(i)*set.dx
	}
	return result
}

// N creates a new set with a new number of values
func (set Set) N(n int) (Set, error) {
	return NewSetN(set.xmin, set.xmax, n)
}

// Dx creates a new set by changing the dx value
func (set Set) Dx(dx float64) (Set, error) {
	return NewSet(set.xmin, set.xmax, dx)
}
