package builder

import (
	"fugologic/fuzzy"
)

// Builder groups custom connector and implication
type Builder struct {
	cnt    fuzzy.Operator
	impl   fuzzy.Implication
	agg    fuzzy.Aggregation
	defuzz fuzzy.Defuzzification

	rules []fuzzy.Rule
}

// NewBuilder creates a builder with a default configuration
func NewBuilder(
	cnt fuzzy.Operator,
	impl fuzzy.Implication,
	agg fuzzy.Aggregation,
	defuzz fuzzy.Defuzzification,
) Builder {
	return Builder{
		cnt:    cnt,
		impl:   impl,
		agg:    agg,
		defuzz: defuzz,
	}
}

// If starts a rule expression
func (bld *Builder) If(premise fuzzy.Premise) expression {
	return expression{
		bld:   bld,
		fzExp: fuzzy.NewExpression([]fuzzy.Premise{premise}, nil),
	}
}

// Engine created using the defined rules and the default configuration
func (bld Builder) Engine() (fuzzy.Engine, error) {
	return fuzzy.NewEngine(bld.rules, bld.agg, bld.defuzz)
}

// add a new rule to the builder
func (bld *Builder) add(rule fuzzy.Rule) {
	bld.rules = append(bld.rules, rule)
}
