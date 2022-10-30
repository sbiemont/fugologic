package fuzzy

import (
	"errors"
	"math"
)

// Premise is an Expression or an IDSet
// A premise can be evaluated or linked to another premise
type Premise interface {
	Evaluate(input DataInput) (float64, error)
}

// Connector links a list of premises
type Connector func(a, b float64) float64

// Operator defines the connectors for a predefined family
// https://commons.wikimedia.org/wiki/Fuzzy_operator
type Operator struct {
	And  Connector
	Or   Connector
	XOr  Connector
	NAnd Connector
	NOr  Connector
}

// OperatorZadeh defines a list of Zadeh connectors
var OperatorZadeh = Operator{
	// Zadeh AND = min
	And: math.Min,
	// Zadeh OR = max
	Or: math.Max,
	// Zadeh XOR = a+b-2*min(a,b)
	XOr: func(a, b float64) float64 { return a + b - 2*math.Min(a, b) },
	// Zadeh NAND = 1-AND = 1-min(a,b)
	NAnd: func(a, b float64) float64 { return 1 - math.Min(a, b) },
	// Zadeh NOR = 1-OR = 1-max(a,b)
	NOr: func(a, b float64) float64 { return 1 - math.Max(a, b) },
}

// OperatorHyperbolic defines a list of hyperbolic connectors
var OperatorHyperbolic = Operator{
	// Hyperbolic AND = a*b
	And: func(a, b float64) float64 { return a * b },
	// Hyperbolic OR = a+b-a*b
	Or: func(a, b float64) float64 { return a + b - a*b },
	// Hyperbolic XOR = a+b-2*a*b
	XOr: func(a, b float64) float64 { return a + b - 2*a*b },
	// Hyperbolic NAND = 1-AND = 1-a*b
	NAnd: func(a, b float64) float64 { return 1 - a*b },
	// Hyperbolic NOR = 1-OR = 1-a-b+a*b
	NOr: func(a, b float64) float64 { return 1 - a - b + a*b },
}

// Expression connects a list of premises. Eg.: A or B or C
// Eg.:
//   - Expression1 = A or B or C
//   - Expression2 = D or E
//   - Expression3 = Expression1 and Expression2 = (A or B or C) and (D or E)
type Expression struct {
	premises []Premise
	connect  Connector
}

// NewExpression initialise a fully evaluable expression
func NewExpression(premises []Premise, connect Connector) Expression {
	return Expression{
		premises: premises,
		connect:  connect,
	}
}

// Connect the current expression, using the connector and the given premise
// Returns <new exp> = <exp> <connect> <premise>
// E.g:    <A and B> = <A>   <and>     <B>
func (exp Expression) Connect(premise Premise, connect Connector) Expression {
	if exp.connect == nil {
		// Direct connection
		return NewExpression(append(exp.premises, premise), connect)
	}

	// Connect both premises in a new expression
	return NewExpression([]Premise{exp, premise}, connect)
}

// Evaluate the expression content
func (exp Expression) Evaluate(input DataInput) (float64, error) {
	// Check
	if len(exp.premises) == 0 {
		return 0, errors.New("expression: at least 1 premise expected")
	}

	// Evaluate premises to compute values
	values := make([]float64, len(exp.premises))
	for i, premise := range exp.premises {
		value, err := premise.Evaluate(input)
		if err != nil {
			return 0, err
		}
		values[i] = value
	}

	// Connect values
	y := values[0]
	if exp.connect != nil {
		for _, value := range values[1:] {
			y = exp.connect(y, value)
		}
	}

	return y, nil
}
