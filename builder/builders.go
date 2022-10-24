package builder

import "fugologic/fuzzy"

// NewFuzzyLogicMamdani sets the default configuration for Mamdani inference system
func NewFuzzyLogicMamdani() FuzzyLogic {
	return NewFuzzyLogic(
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
