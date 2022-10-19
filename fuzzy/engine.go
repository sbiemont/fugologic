package fuzzy

import (
	"fmt"
	"fugologic/id"

	"golang.org/x/sync/errgroup"
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

// Evalute rules (in parallel) and defuzz result
func (eng Engine) Evaluate(input DataInput) (DataOutput, error) {
	evaluatedIDSets := make([][]IDSet, len(eng.rules)) // prepare results for go routines
	var grp errgroup.Group
	for i, rule := range eng.rules {
		iCpy := i
		ruleCpy := rule
		grp.Go(func() error {
			var errEval error
			evaluatedIDSets[iCpy], errEval = ruleCpy.evaluate(input)
			return errEval
		})
	}

	// Wait for all evaluations
	err := grp.Wait()
	if err != nil {
		return nil, err
	}

	// Push result into the defuzzer
	var flattenIDSets []IDSet
	for _, idSets := range evaluatedIDSets {
		flattenIDSets = append(flattenIDSets, idSets...)
	}

	// Apply defuzzification
	dfz := newDefuzzer(eng.defuzz, eng.agg)
	return dfz.defuzz(flattenIDSets), nil
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
	idVals := IDSets(idSets).extractIDVal()

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
