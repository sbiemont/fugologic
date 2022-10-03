package fuzzy

import (
	"fmt"
	"fugologic/id"
	"math"
)

var (
	GAUSS    = "gauss"
	GBELL    = "gbell"
	TRAP     = "trap"
	TRI      = "tri"
	STEPUP   = "step-up"
	STEPDOWN = "step-down"
)

// SetBuilder helps create a new set
type SetBuilder interface {
	New() (Set, error)
}

// helper: checks that the input parameters are sorted
func checkSorted(name string, params ...float64) error {
	if len(params) == 0 {
		return nil
	}
	a := params[0]
	for _, p := range params[1:] {
		if a > p {
			return fmt.Errorf("%s: params shall be sorted", name)
		}
		a = p
	}
	return nil
}

// Gauss builder
type Gauss struct {
	Sigma, C float64
}

// New Gauss membership function
// https://www.mathworks.com/help/fuzzy/gaussmf.html
//
// _/¯\_
func (set Gauss) New() (Set, error) {
	if set.Sigma == 0 {
		return nil, fmt.Errorf("%s: first parameter must be non zero", GAUSS)
	}
	sigma2 := 2.0 * math.Pow(set.Sigma, 2.0)
	return func(x float64) float64 {
		return float64(math.Exp(-(math.Pow(x-set.C, 2.0)) / sigma2))
	}, nil
}

// Gbell builder
type Gbell struct {
	A, B, C float64
}

// New generalized bell-shaped membership function: 1 / (1 + ((x-c)/a)^2b)
// https://www.mathworks.com/help/fuzzy/gbellmf.html
//
// _/¯\_
func (set Gbell) New() (Set, error) {
	if set.A == 0 {
		return nil, fmt.Errorf("%s: first parameter must be non zero", GBELL)
	}
	b2 := set.B * 2
	return func(x float64) float64 {
		return 1.0 / (1.0 + math.Pow(math.Abs((x-set.C)/set.A), b2))
	}, nil
}

// Trapezoid builder
type Trapezoid struct {
	A, B, C, D float64
}

// New trapezoid membership function
// https://www.mathworks.com/help/fuzzy/trapmf.html
// - a: first base (left to right) of the function (y=0)
// - b: peak of the function (y=1)
// - c: second peak of the function (y=1)
// - d: last base of the function (y=0)
//
// _/¯\_
func (set Trapezoid) New() (Set, error) {
	if err := checkSorted(TRAP, set.A, set.B, set.C, set.D); err != nil {
		return nil, err
	}

	ba := set.B - set.A
	dc := set.D - set.C
	return func(x float64) float64 {
		switch {
		case set.A < x && x <= set.B:
			return (x - set.A) / ba
		case set.B <= x && x <= set.C:
			return 1
		case set.C <= x && x < set.D:
			return (set.D - x) / dc
		default:
			return 0
		}
	}, nil
}

// Triangular builder
type Triangular struct {
	A, B, C float64
}

// New triangular membership function
// https://www.mathworks.com/help/fuzzy/trimf.html
// - a: first base (left to right) of the function (y=0)
// - b: peak of the function (y=1)
// - c: second base of the function (y=0)
//
// _/\_
func (set Triangular) New() (Set, error) {
	if err := checkSorted(TRI, set.A, set.B, set.C); err != nil {
		return nil, err
	}

	ba := set.B - set.A
	cb := set.C - set.B
	return func(x float64) float64 {
		switch {
		case set.A < x && x <= set.B:
			return (x - set.A) / ba
		case set.B <= x && x < set.C:
			return (set.C - x) / cb
		default:
			return 0
		}
	}, nil
}

// StepUp builder
type StepUp struct {
	A, B float64
}

// New step-up membership function
// - a: first base (left to right) of the function (y=0)
// - b: peak of the function (y=1)
//
// _/¯
func (set StepUp) New() (Set, error) {
	if err := checkSorted(STEPUP, set.A, set.B); err != nil {
		return nil, err
	}

	ab := set.A - set.B
	return func(x float64) float64 {
		switch {
		case x >= set.B:
			return 1
		case set.A < x && x <= set.B:
			return (set.A - x) / ab
		default:
			return 0
		}
	}, nil
}

// StepDown builder
type StepDown struct {
	A, B float64
}

// NewSetStepDown membership function
// - a: peak (left to right) of the function (1)
// - b: last base of the function (0)
//
// ¯\_
func (set StepDown) New() (Set, error) {
	if err := checkSorted(STEPDOWN, set.A, set.B); err != nil {
		return nil, err
	}

	ba := set.B - set.A
	return func(x float64) float64 {
		switch {
		case x <= set.A:
			return 1
		case set.A <= x && x < set.B:
			return (set.B - x) / ba
		default:
			return 0
		}
	}, nil
}

// NewIDSets builds a list of named fuzzy sets
func NewIDSets(fsets map[id.ID]SetBuilder) (map[id.ID]Set, error) {
	sets := make(map[id.ID]Set, len(fsets))
	for uuid, fset := range fsets {
		set, err := fset.New()
		if err != nil {
			return nil, fmt.Errorf("%s: %w", uuid, err)
		}
		sets[uuid] = set
	}
	return sets, nil
}
