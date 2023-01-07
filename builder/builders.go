package builder

import "github.com/sbiemont/fugologic/fuzzy"

// Mamdani predefined configuration
func Mamdani() Config {
	return Config{
		Optr:   fuzzy.OperatorZadeh{},
		Impl:   fuzzy.ImplicationMin,
		Agg:    fuzzy.AggregationUnion,
		Defuzz: fuzzy.DefuzzificationCentroid,
	}
}

// Config gathers the configuration for a fuzzy rule builder
type Config struct {
	Optr   fuzzy.Operator
	Impl   fuzzy.Implication
	Agg    fuzzy.Aggregation
	Defuzz fuzzy.Defuzzification
}

// FuzzyLogic returns a fuzzy-logic rules builder using the current configuration
func (cfg Config) FuzzyLogic() FuzzyLogic {
	return NewFuzzyLogic(
		cfg.Optr,
		cfg.Impl,
		cfg.Agg,
		cfg.Defuzz,
	)
}

// FuzzyAssoMatrix returns a fuzzy-associative-matrix builder using the current configuration
func (cfg Config) FuzzyAssoMatrix() FuzzyAssoMatrix {
	return NewFuzzyAssoMatrix(
		cfg.Optr,
		cfg.Impl,
		cfg.Agg,
		cfg.Defuzz,
	)
}
