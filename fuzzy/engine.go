package fuzzy

// Engine is responsible for evaluating all rules and defuzzing
type Engine struct {
	rules  []Rule
	defuzz Defuzzer
}

// NewEngine builds a new Engine instance
func NewEngine(rules []Rule, defuzz Defuzzer) (Engine, error) {
	// Gather inputs and outpus
	var idSets []IDSet
	for _, rule := range rules {
		idSets = append(idSets, rule.Inputs()...)
		idSets = append(idSets, rule.outputs...)
	}

	// Check
	if err := newChecker(idSets).check(); err != nil {
		return Engine{}, err
	}

	return Engine{
		rules:  rules,
		defuzz: defuzz,
	}, nil
}

// Evalute rules and defuzz result
func (eng Engine) Evaluate(input DataInput) (DataOutput, error) {
	for _, rule := range eng.rules {
		// Evaluate rule
		idSets, err := rule.evaluate(input)
		if err != nil {
			return nil, err
		}

		// Push result into the defuzzer
		eng.defuzz.Add(idSets)
	}

	// Apply defuzzification
	return eng.defuzz.Defuzz()
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
