package fuzzy

import (
	"math"
)

// Set defines a Fuzzy Set Type-I
// https://en.wikipedia.org/wiki/Fuzzy_set
type Set func(float64) float64

// chain 2 functions as one
func (fs Set) chain(fs2 Set, merge func(float64, float64) float64) Set {
	return func(x float64) float64 {
		return merge(fs(x), fs2(x))
	}
}

// Union of 2 sets
func (fs Set) Union(fs2 Set) Set {
	return fs.chain(fs2, math.Max)
}

// Intersection of 2 sets
func (fs Set) Intersection(fs2 Set) Set {
	return fs.chain(fs2, math.Min)
}

// Complement of the current set
func (fs Set) Complement() Set {
	return func(x float64) float64 {
		return 1.0 - fs(x)
	}
}

// Min merges a minimum membership method
func (fs Set) Min(y float64) Set {
	return func(x float64) float64 {
		return math.Min(fs(x), y)
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
		return math.Max(math.Min(math.Min(((x-a)/(b-a)), 1), (d-x)/(d-c)), 0)
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
		return math.Min(math.Max((a-x)/(a-b), 0), 1)
	}
}

// NewSetStepDown membership function
//  * a: peak (left to right) of the function (1)
//  * b: last base of the function (0)
// _
//  \_
func NewSetStepDown(a, b float64) Set {
	return func(x float64) float64 {
		return math.Min(math.Max((b-x)/(b-a), 0), 1)
	}
}
