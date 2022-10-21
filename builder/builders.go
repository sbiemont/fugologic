package builder

import "fugologic/fuzzy"

// NewBuilderMamdani sets the default configuration for Mamdani inference system
func NewBuilderMamdani() Builder {
	return NewBuilder(
		fuzzy.OperatorZadeh,
		fuzzy.ImplicationMin,
		fuzzy.AggregationUnion,
		fuzzy.DefuzzificationCentroid,
	)
}

// NewFuzzyAssoMatrixMamdani sets the default configuration for a fuzzy associative matrix builder
func NewFuzzyAssoMatrixMamdani() FuzzyAssoMatrix {
	return NewFuzzyAssoMatrix(
		fuzzy.OperatorZadeh.And,
		fuzzy.ImplicationMin,
		fuzzy.AggregationUnion,
		fuzzy.DefuzzificationCentroid,
	)
}
