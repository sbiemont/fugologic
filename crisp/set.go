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

// Values translates the interval into discrete increasing values
func (set Set) Values() []float64 {
	if set.dx == 0 {
		return nil
	}

	result := []float64{}
	for i := set.xmin; i <= set.xmax; i += set.dx {
		result = append(result, i)
	}
	return result
}
