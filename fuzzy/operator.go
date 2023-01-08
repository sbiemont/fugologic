package fuzzy

import "math"

// Operator defines the connectors for a predefined family
// https://commons.wikimedia.org/wiki/Fuzzy_operator
// Functions are connectors (see. Connector)
type Operator interface {
	And(a, b float64) float64
	Or(a, b float64) float64
	XOr(a, b float64) float64
}

// OperatorZadeh defines a list of Zadeh connectors
type OperatorZadeh struct{}

func (OperatorZadeh) And(a, b float64) float64 { return math.Min(a, b) }
func (OperatorZadeh) Or(a, b float64) float64  { return math.Max(a, b) }
func (OperatorZadeh) XOr(a, b float64) float64 { return a + b - 2*math.Min(a, b) }

// OperatorHyperbolic defines a list of hyperbolic connectors
type OperatorHyperbolic struct{}

func (OperatorHyperbolic) And(a, b float64) float64 { return a * b }
func (OperatorHyperbolic) Or(a, b float64) float64  { return a + b - a*b }
func (OperatorHyperbolic) XOr(a, b float64) float64 { return a + b - 2*a*b }
