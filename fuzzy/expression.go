package fuzzy

import (
	"errors"
)

// Premise is an Expression or an IDSet
// A premise can be evaluated or linked to another premise
type Premise interface {
	Evaluate(input DataInput) (float64, error)
}

// Connector links a list of premises
type Connector func(a, b float64) float64

// Expression connects a list of premises. Eg.: A or B or C
// Eg.:
//   - Expression1 = A or B or C
//   - Expression2 = D or E
//   - Expression3 = Expression1 and Expression2 = (A or B or C) and (D or E)
type Expression struct {
	premises   []Premise // List all premises to be connected
	connect    Connector // Connector to be applied on the premises
	complement bool      // Complement (false by default)
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

// Not complements the current expression
func (exp Expression) Not() Expression {
	return Expression{
		premises:   exp.premises,
		connect:    exp.connect,
		complement: !exp.complement,
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

	// Apply complement
	if exp.complement {
		y = 1 - y
	}

	return y, nil
}
