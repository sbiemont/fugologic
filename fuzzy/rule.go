package fuzzy

import (
	"errors"
	"math"
)

// Connector links a list of premises
type Connector func(a, b float64) float64

var (
	ConnectorNone Connector = nil
	ConnectorAnd  Connector = math.Min
	ConnectorOr   Connector = math.Max
)

// Premise can be evaluated (like a fuzzy set or an expression)
type Premise interface {
	Evaluate(input DataInput) (float64, error)
}

// Expression describes "connect(premises)". Eg.: A or B or C
// An Expression is also a Premise.
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
	if len(exp.premises) < 1 {
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

// flatten expression into a list of IDSet
func (exp Expression) flatten() []IDSet {
	return exp.extractSets(nil, []Premise{exp})
}

// extractSets extracts the IDSets from a list of premises
func (exp Expression) extractSets(idSets []IDSet, premises []Premise) []IDSet {
	for _, premise := range premises {
		switch p := premise.(type) {
		case IDSet:
			idSets = append(idSets, p)
		case Expression:
			idSets = exp.extractSets(idSets, p.premises)
		}
	}
	return idSets
}

type Implication func(set Set, y float64) Set

var (
	// ImplicationProd returns the product of a Set
	ImplicationProd Implication = func(set Set, y float64) Set { return set.Multiply(y) }

	// ImplicationMin sets the max upper bound
	ImplicationMin Implication = func(set Set, y float64) Set { return set.Min(y) }
)

// Rule evaluates the input expression + implication + fuzzy output
type Rule struct {
	inputs      Expression
	implication Implication
	outputs     []IDSet
}

// NewRule builds a new Rule instance
func NewRule(inputs Expression, implication Implication, outputs []IDSet) Rule {
	return Rule{
		inputs:      inputs,
		implication: implication,
		outputs:     outputs,
	}
}

// evaluate and return the fuzzy output using crisp input
// Outputs
// * One fuzzy Set for each output
// * An error is returned if inputs are missing
func (rule Rule) evaluate(input DataInput) ([]IDSet, error) {
	// Evaluate inputs
	y, err := rule.inputs.Evaluate(input)
	if err != nil {
		return nil, err
	}

	// Evaluate outputs => create a NEW fuzzy Set with the same output ID
	result := make([]IDSet, len(rule.outputs))
	for i, out := range rule.outputs {
		result[i] = IDSet{
			uuid:   out.uuid,
			set:    rule.implication(out.set, y),
			parent: out.parent,
		}
	}

	return result, nil
}

// Inputs gather and flatten all IDSet from rules' expressions
func (rule Rule) Inputs() []IDSet {
	return rule.inputs.flatten()
}
