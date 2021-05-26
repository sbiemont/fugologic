package fuzzy

import (
	"math"
)

// Set defines a Fuzzy Set Type-I
// https://en.wikipedia.org/wiki/Fuzzy_set
type Set func(float64) float64

var (
	union        = math.Max
	intersection = math.Min
)

// aggregate of 2 sets using a specific method
// E.g: fs1.aggregate(fs2, math.Max)
func (fs Set) aggregate(fs2 Set, fct func(float64, float64) float64) Set {
	return func(x float64) float64 {
		return fct(fs(x), fs2(x))
	}
}

// Union of 2 sets
func (fs Set) Union(fs2 Set) Set {
	return fs.aggregate(fs2, union)
}

// Intersection of 2 sets
func (fs Set) Intersection(fs2 Set) Set {
	return fs.aggregate(fs2, intersection)
}

// Complement of the current set
func (fs Set) Complement() Set {
	return func(x float64) float64 {
		return 1.0 - fs(x)
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

// NewSetGauss membership function
// https://www.mathworks.com/help/fuzzy/gaussmf.html
//   _
// _/ \_
func NewSetGauss(sigma, c float64) Set {
	return func(x float64) float64 {
		return float64(math.Exp(-(math.Pow(x-c, 2.0)) / (2.0 * math.Pow(sigma, 2.0))))
	}
}

// NewSetGbell membership function: 1 / (1 + ((x-c)/a)^2b)
// https://www.mathworks.com/help/fuzzy/gbellmf.html
//   _
// _/ \_
func NewSetGbell(a, b, c float64) Set {
	return func(x float64) float64 {
		return 1.0 / (1.0 + math.Pow(math.Abs((x-c)/a), 2*b))
	}
}

// NewSetTrapezoid membership function
// https://www.mathworks.com/help/fuzzy/trapmf.html
//  * a: first base (left to right) of the function (y=0)
//  * b: peak of the function (y=1)
//  * c: seconde base of the function (y=0)
//   _
// _/ \_
func NewSetTrapezoid(a, b, c, d float64) Set {
	return func(x float64) float64 {
		switch {
		case a < x && x <= b:
			return (x - a) / (b - a)
		case b <= x && x <= c:
			return 1
		case c <= x && x < d:
			return (d - x) / (d - c)
		default:
			return 0
		}
	}
}

// NewSetTriangular membership function
// https://www.mathworks.com/help/fuzzy/trimf.html
//  * a: first base (left to right) of the function (y=0)
//  * b: peak of the function (y=1)
//  * c: seconde base of the function (y=0)
// _/\_
func NewSetTriangular(a, b, c float64) Set {
	return func(x float64) float64 {
		switch {
		case a < x && x <= b:
			return (x - a) / (b - a)
		case b <= x && x < c:
			return (c - x) / (c - b)
		default:
			return 0
		}
	}
}

// NewSetStepUp membership function
//  * a: first base (left to right) of the function (y=0)
//  * b: peak of the function (y=1)
//   _
// _/
func NewSetStepUp(a, b float64) Set {
	return func(x float64) float64 {
		switch {
		case x >= b:
			return 1
		case a < x && x <= b:
			return (a - x) / (a - b)
		default:
			return 0
		}
	}
}

// NewSetStepDown membership function
//  * a: peak (left to right) of the function (1)
//  * b: last base of the function (0)
// _
//  \_
func NewSetStepDown(a, b float64) Set {
	return func(x float64) float64 {
		switch {
		case x <= a:
			return 1
		case a <= x && x < b:
			return (b - x) / (b - a)
		default:
			return 0
		}
	}
}
