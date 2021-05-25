package builder

import "fugologic.git/fuzzy"

// Builder groups custom connector and implication
type Builder struct {
	and  fuzzy.Connector
	or   fuzzy.Connector
	impl fuzzy.Implication
}

// NewBuilder sets the custom operations and creates a new Builder instance
func NewBuilder(and fuzzy.Connector, or fuzzy.Connector, impl fuzzy.Implication) Builder {
	return Builder{
		and:  and,
		or:   or,
		impl: impl,
	}
}

// If starts a rule expression
func (bld Builder) If(premise fuzzy.Premise) Expression {
	return Expression{
		builder: bld,
		fzExp:   fuzzy.NewExpression([]fuzzy.Premise{premise}, nil),
	}
}
