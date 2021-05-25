package builder

import "fugologic.git/fuzzy"

// Expression embeds a custom builder and a fuzzy expression
type Expression struct {
	builder Builder
	fzExp   fuzzy.Expression
}

// Evaluate the fuzzy expression linked
func (exp Expression) Evaluate(input fuzzy.DataInput) (float64, error) {
	return exp.fzExp.Evaluate(input)
}

// And connects the current expression and a premise with the AND connector of the builder
func (exp Expression) And(premise fuzzy.Premise) Expression {
	return Expression{
		builder: exp.builder,
		fzExp:   exp.fzExp.Connect(premise, exp.builder.and),
	}
}

// Or connects the current expression and a premise with the OR connector of the builder
func (exp Expression) Or(premise fuzzy.Premise) Expression {
	return Expression{
		builder: exp.builder,
		fzExp:   exp.fzExp.Connect(premise, exp.builder.or),
	}
}

// Then describes the consequence of an implication
func (exp Expression) Then(consequence []fuzzy.IDSet) fuzzy.Rule {
	return fuzzy.NewRule(
		exp.fzExp,
		exp.builder.impl,
		consequence,
	)
}
