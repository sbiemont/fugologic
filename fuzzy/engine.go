package fuzzy

import (
	"fmt"
	"fugologic/id"
)

// Engine is responsible for evaluating all rules and defuzzing
type Engine struct {
	uuid   id.ID // optional
	rules  []Rule
	agg    Aggregation
	defuzz Defuzzification
}

// NewEngine builds a new Engine instance
//   - The Aggregation merges all result sets together
//   - The Defuzzification extracts one value from the aggregation
func NewEngine(r []Rule, agg Aggregation, defuzz Defuzzification) (Engine, error) {
	// Check
	inputs, outputs := rules(r).io()
	if err := checkIDs(append(inputs, outputs...)); err != nil {
		return Engine{}, err
	}

	return Engine{
		rules:  r,
		defuzz: defuzz,
		agg:    agg,
	}, nil
}

// Evalute rules and defuzz result
func (eng Engine) Evaluate(input DataInput) (DataOutput, error) {
	var evaluatedIDSets []IDSet
	for _, rule := range eng.rules {
		// Evaluate rule
		idSets, err := rule.evaluate(input)
		if err != nil {
			return nil, err
		}

		// Push result into the defuzzer
		evaluatedIDSets = append(evaluatedIDSets, idSets...)
	}

	// Apply defuzzification
	dfz := newDefuzzer(eng.defuzz, eng.agg)
	return dfz.defuzz(evaluatedIDSets), nil
}

// FlattenIO gather and flatten all IDSet from rules' expressions
// Return inputs and outputs IDSet
func (eng Engine) io() ([]IDSet, []IDSet) {
	return rules(eng.rules).io()
}

// checkIDs of a list of IDSet
// Get all unique IDVal, check them and their whole IDSet
func checkIDs(idSets []IDSet) error {
	// Extract all unique IDVal
	var idVals = make(map[*IDVal]struct{})
	for _, idSet := range idSets {
		idVals[idSet.parent] = struct{}{}
	}

	// Extract all ids of IDVals and compare them
	uniqueIDs := make(map[id.ID]struct{})
	for idVal := range idVals {
		uuid := idVal.ID()
		if _, exists := uniqueIDs[uuid]; exists {
			return fmt.Errorf("values: id `%s` already defined", uuid)
		}
		uniqueIDs[uuid] = struct{}{}
	}

	return nil
}
