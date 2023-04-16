package fuzzy

import (
	"math"
)

// Set defines a Fuzzy Set Type-1
// https://en.wikipedia.org/wiki/Fuzzy_set
type Set func(float64) float64

// aggregate of 2 sets using a specific method
// E.g: fs1.aggregate(fs2, math.Max)
func (fs Set) aggregate(fs2 Set, fct func(float64, float64) float64) Set {
	return func(x float64) float64 {
		return fct(fs(x), fs2(x))
	}
}

// Min merges a minimum membership method
func (fs Set) Min(k float64) Set {
	return func(x float64) float64 {
		return math.Min(fs(x), k)
	}
}

// Multiply the current membership method with a constant factor
func (fs Set) Multiply(k float64) Set {
	return func(x float64) float64 {
		return fs(x) * k
	}
}
