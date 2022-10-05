package builder

import "fugologic/fuzzy"

// NewBuilderMamdani sets the default configuration for Mamdani inference system
func NewBuilderMamdani() Builder {
	return NewBuilder(
		Connector{
			And:  fuzzy.ConnectorZadehAnd,
			Or:   fuzzy.ConnectorZadehOr,
			XOr:  fuzzy.ConnectorZadehXOr,
			NAnd: fuzzy.ConnectorZadehNAnd,
			NOr:  fuzzy.ConnectorZadehNOr,
		},
		fuzzy.ImplicationMin,
		fuzzy.AggregationUnion,
		fuzzy.DefuzzificationCentroid,
	)
}
