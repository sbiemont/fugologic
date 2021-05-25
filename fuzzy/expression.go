package fuzzy

import (
	"errors"
	"math"
)

// Premise is an Expression or an IDSet
// A premise can be evaluated or linked to another premise
type Premise interface {
	Evaluate(input DataInput) (float64, error)
	And(premise Premise) Premise // use default `And` connector
	Or(premise Premise) Premise  // use default `Or` connector
}

// If starts a rule expression
func If(premise Premise) RuleExpression {
	return RuleExpression{
		premise: premise,
	}
}

// Connector links a list of premises
type Connector func(a, b float64) float64

// https://commons.wikimedia.org/wiki/Fuzzy_operator
var (
	// Zadeh AND = min
	ConnectorZadehAnd Connector = math.Min
	// Zadeh OR = max
	ConnectorZadehOr Connector = math.Max
	// Zadeh XOR = a+b-2*min(a,b)
	ConnectorZadehXor Connector = func(a, b float64) float64 { return a + b - 2*math.Min(a, b) }
	// Zadeh NAND = 1-max(a,b)
	ConnectorZadehNand Connector = func(a, b float64) float64 { return 1 - math.Min(a, b) }
	// Zadeh NOR = 1-max(a,b)
	ConnectorZadehNor Connector = func(a, b float64) float64 { return 1 - math.Max(a, b) }

	// Hyperbolic AND = a*b
	ConnectorHyperbolicAnd Connector = func(a, b float64) float64 { return a * b }
	// Hyperbolic OR = a+b-a*b
	ConnectorHyperbolicOr Connector = func(a, b float64) float64 { return a + b - a*b }
	// Hyperbolic XOR = a+b-2*a*b
	ConnectorHyperbolicXor Connector = func(a, b float64) float64 { return a + b - 2*a*b }
	// Hyperbolic NAND = 1-a*b
	ConnectorHyperbolicNand Connector = func(a, b float64) float64 { return 1 - a*b }
	// Hyperbolic OR = 1-a-b+a*b
	ConnectorHyperbolicNor Connector = func(a, b float64) float64 { return 1 - a - b + a*b }
)

// Expression connects a list of premises. Eg.: A or B or C
// Eg.:
//  * Expression1 = A or B or C
//  * Expression2 = D or E
//  * Expression3 = Expression1 and Expression2 = (A or B or C) and (D or E)
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

// aggregate the given premise to the current expression using a custom connector
func (exp Expression) aggregate(premise Premise, connect Connector) Premise {
	if exp.connect == nil {
		// Direct connection
		return Expression{
			premises: append(exp.premises, premise),
			connect:  connect,
		}
	}

	// Connect both premises in a new expression
	return Expression{
		premises: []Premise{exp, premise},
		connect:  connect,
	}
}

// And connects the current expression and a premise with the default AND connector
func (exp Expression) And(premise Premise) Premise {
	return exp.aggregate(premise, ConnectorZadehAnd)
}

// Or connects the current expression and a premise with the default OR connector
func (exp Expression) Or(premise Premise) Premise {
	return exp.aggregate(premise, ConnectorZadehOr)
}

// RuleExpression connects an expression to an implication
type RuleExpression struct {
	premise Premise
}

// Then describes the consequence of an implication
func (rexp RuleExpression) Then(consequence []IDSet) Rule {
	return Rule{
		inputs:      rexp.premise,
		implication: ImplicationMin,
		outputs:     consequence,
	}
}
