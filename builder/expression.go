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

// And connects the current expression and a premise with the AND connector of the builder
func (exp expression) And(premise fuzzy.Premise) expression {
	return expression{
		bld:   exp.bld,
		fzExp: exp.fzExp.Connect(premise, exp.bld.and),
	}
}

// Or connects the current expression and a premise with the OR connector of the builder
func (exp expression) Or(premise fuzzy.Premise) expression {
	return expression{
		bld:   exp.bld,
		fzExp: exp.fzExp.Connect(premise, exp.bld.or),
	}
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
