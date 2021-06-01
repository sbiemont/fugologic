package fuzzy

// Engine is responsible for evaluating all rules and defuzzing
type Engine struct {
	rules  []Rule
	agg    Aggregation
	defuzz Defuzzification
}

// NewEngine builds a new Engine instance
//  * The Aggregation merges all result sets together
//  * The Defuzzification extracts one value from the aggregation
func NewEngine(rules []Rule, agg Aggregation, defuzz Defuzzification) (Engine, error) {
	// Gather inputs and outpus
	var idSets []IDSet
	for _, rule := range rules {
		idSets = append(idSets, rule.Inputs()...)
		idSets = append(idSets, rule.outputs...)
	}

	// Check
	if err := checkIDs(idSets); err != nil {
		return Engine{}, err
	}

	return Engine{
		rules:  rules,
		defuzz: defuzz,
		agg:    agg,
	}, nil
}

// Evalute rules and defuzz result
func (eng Engine) Evaluate(input DataInput) (DataOutput, error) {
	dfz := newDefuzzer(eng.defuzz, eng.agg)

	for _, rule := range eng.rules {
		// Evaluate rule
		idSets, err := rule.evaluate(input)
		if err != nil {
			return nil, err
		}

		// Push result into the defuzzer
		dfz.add(idSets)
	}

	// Apply defuzzification
	return dfz.defuzz()
}

// Inputs gather and flatten all IDSet from rules' expressions
// An IDSet is returned only once
func (eng Engine) Inputs() []IDSet {
	var result []IDSet
	for _, rule := range eng.rules {
		result = append(result, rule.Inputs()...)
	}
	return result
}

func (eng Engine) Outputs() []IDSet {
	var result []IDSet
	for _, rule := range eng.rules {
		result = append(result, rule.outputs...)
	}
	return result
}
