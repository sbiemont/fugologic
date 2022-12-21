package builder

import (
	"github.com/sbiemont/fugologic/fuzzy"
)

// FuzzyLogic groups custom connector and implication
type FuzzyLogic struct {
	cnt    fuzzy.Operator
	impl   fuzzy.Implication
	agg    fuzzy.Aggregation
	defuzz fuzzy.Defuzzification

	rules []fuzzy.Rule
}

// NewFuzzyLogic creates a builder with a default configuration
func NewFuzzyLogic(
	cnt fuzzy.Operator,
	impl fuzzy.Implication,
	agg fuzzy.Aggregation,
	defuzz fuzzy.Defuzzification,
) FuzzyLogic {
	return FuzzyLogic{
		cnt:    cnt,
		impl:   impl,
		agg:    agg,
		defuzz: defuzz,
	}
}

// If starts a rule expression
func (fl *FuzzyLogic) If(premise fuzzy.Premise) flExpression {
	return flExpression{
		fl:    fl,
		fzExp: fuzzy.NewExpression([]fuzzy.Premise{premise}, nil),
	}
}

// Engine created using the defined rules and the default configuration
func (fl FuzzyLogic) Engine() (fuzzy.Engine, error) {
	return fuzzy.NewEngine(fl.rules, fl.agg, fl.defuzz)
}

// add a new rule to the builder
func (fl *FuzzyLogic) add(rule fuzzy.Rule) {
	fl.rules = append(fl.rules, rule)
}

// flExpression embeds a custom builder and a fuzzy flExpression
type flExpression struct {
	fl    *FuzzyLogic
	fzExp fuzzy.Expression
}

// Evaluate the fuzzy expression linked
func (exp flExpression) Evaluate(input fuzzy.DataInput) (float64, error) {
	return exp.fzExp.Evaluate(input)
}

// connect the current expression with a new one with a connector
func (exp flExpression) connect(premise fuzzy.Premise, cnt fuzzy.Connector) flExpression {
	return flExpression{
		fl:    exp.fl,
		fzExp: exp.fzExp.Connect(premise, cnt),
	}
}

// And connects the current expression and a premise with the AND connector of the builder
func (exp flExpression) And(premise fuzzy.Premise) flExpression {
	return exp.connect(premise, exp.fl.cnt.And)
}

// Or connects the current expression and a premise with the OR connector of the builder
func (exp flExpression) Or(premise fuzzy.Premise) flExpression {
	return exp.connect(premise, exp.fl.cnt.Or)
}

// NOr connects the current expression and a premise with the NOT-OR connector of the builder
func (exp flExpression) NOr(premise fuzzy.Premise) flExpression {
	return exp.connect(premise, exp.fl.cnt.NOr)
}

// NAnd connects the current expression and a premise with the NOT-AND connector of the builder
func (exp flExpression) NAnd(premise fuzzy.Premise) flExpression {
	return exp.connect(premise, exp.fl.cnt.NAnd)
}

// XOr connects the current expression and a premise with the XOR connector of the builder
func (exp flExpression) XOr(premise fuzzy.Premise) flExpression {
	return exp.connect(premise, exp.fl.cnt.XOr)
}

// Then describes the consequence of an implication AND stores the rule into the builder
// At least one consequence is expected
func (exp flExpression) Then(consequence ...fuzzy.IDSet) {
	rule := fuzzy.NewRule(
		exp.fzExp,
		exp.fl.impl,
		consequence,
	)

	// Add the rule to the builder
	exp.fl.add(rule)
}
