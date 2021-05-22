package crisp

// Set (or universe) of possible discrete `x` values
type Set struct {
	xmin float64
	xmax float64
	dx   float64
}

// NewSetDx builds a new Set instance
func NewSetDx(xmin, xmax, dx float64) Set {
	return Set{
		xmin: xmin,
		xmax: xmax,
		dx:   dx,
	}
}

// Values translates the interval into discrete increasing values
func (set Set) Values() []float64 {
	result := []float64{}
	for i := set.xmin; i <= set.xmax; i += set.dx {
		result = append(result, i)
	}
	return result
}
