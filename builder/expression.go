package builder

import "fugologic/fuzzy"

// expression embeds a custom builder and a fuzzy expression
type expression struct {
	bld   *Builder
	fzExp fuzzy.Expression
}

// Evaluate the fuzzy expression linked
func (exp expression) Evaluate(input fuzzy.DataInput) (float64, error) {
	return exp.fzExp.Evaluate(input)
}

// connect the current expression with a new one with a connector
func (exp expression) connect(premise fuzzy.Premise, cnt fuzzy.Connector) expression {
	return expression{
		bld:   exp.bld,
		fzExp: exp.fzExp.Connect(premise, cnt),
	}
}

// And connects the current expression and a premise with the AND connector of the builder
func (exp expression) And(premise fuzzy.Premise) expression {
	return exp.connect(premise, exp.bld.cnt.And)
}

// Or connects the current expression and a premise with the OR connector of the builder
func (exp expression) Or(premise fuzzy.Premise) expression {
	return exp.connect(premise, exp.bld.cnt.Or)
}

// NOr connects the current expression and a premise with the NOT-OR connector of the builder
func (exp expression) NOr(premise fuzzy.Premise) expression {
	return exp.connect(premise, exp.bld.cnt.NOr)
}

// NAnd connects the current expression and a premise with the NOT-AND connector of the builder
func (exp expression) NAnd(premise fuzzy.Premise) expression {
	return exp.connect(premise, exp.bld.cnt.NAnd)
}

// Then describes the consequence of an implication AND stores the rule into the builder
// At least one consequence is expected
func (exp expression) Then(consequence ...fuzzy.IDSet) {
	rule := fuzzy.NewRule(
		exp.fzExp,
		exp.bld.impl,
		consequence,
	)

	// Add the rule to the builder
	exp.bld.add(rule)
}
