package builder

import "fugologic/fuzzy"

// NewBuilderMamdani sets the default configuration for Mamdani inference system
func NewBuilderMamdani() Builder {
	return NewBuilder(
		fuzzy.ConnectorZadehAnd,
		fuzzy.ConnectorZadehOr,
		fuzzy.ImplicationMin,
		fuzzy.AggregationUnion,
		fuzzy.DefuzzificationCentroid,
	)
}

// NewBuilderSugeno sets the default configuration for Takagi-Sugeno inference system
func NewBuilderSugeno() Builder {
	return NewBuilder(
		fuzzy.ConnectorZadehAnd,
		fuzzy.ConnectorZadehOr,
		fuzzy.ImplicationProd,
		fuzzy.AggregationUnion,
		fuzzy.DefuzzificationCentroid,
	)
}
