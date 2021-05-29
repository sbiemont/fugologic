package fuzzy

import (
	"fmt"

	"fugologic.git/id"
)

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

func checkIDs(idSets []IDSet) error {
	// Extract all unique IDVal
	var idVals = make(map[*IDVal]interface{})
	for _, idSet := range idSets {
		idVals[idSet.parent] = nil
	}

	// Extract all IDSets of IDVals and compare them
	var vals []id.Identifiable
	var sets []id.Identifiable
	for idVal := range idVals {
		vals = append(vals, idVal)
		for _, idSet := range idVal.idSets {
			sets = append(sets, idSet)
		}
	}

	// Check all
	if err := id.NewChecker(vals).Check(); err != nil {
		return fmt.Errorf("values: %s", err)
	}
	if err := id.NewChecker(sets).Check(); err != nil {
		return fmt.Errorf("sets: %s", err)
	}
	return nil
}
