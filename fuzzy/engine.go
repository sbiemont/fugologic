package fuzzy

import (
	"fmt"

	"github.com/sbiemont/fugologic/id"

	"golang.org/x/sync/errgroup"
)

// Engine is responsible for evaluating all rules and defuzzing
type Engine struct {
	rules  []Rule
	agg    Aggregation
	defuzz Defuzzification

	maxWorkers int // 0: unlimited, -1: sequential
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
// Use WithMaxWorkers to limit the number of go routines
func (eng Engine) Evaluate(input DataInput) (DataOutput, error) {
	evaluatedIDSets := make([][]IDSet, len(eng.rules)) // prepare results for go routines
	var grp errgroup.Group
	grp.SetLimit(eng.getWorkerLimit(len(eng.rules)))

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

// IO gather and flatten all IDSet from rules' expressions
// Return inputs and outputs IDSet
func (eng Engine) IO() ([]IDSet, []IDSet) {
	return rules(eng.rules).io()
}

// checkIDs of a list of IDSet
// Get all unique IDVal, check them and their whole IDSet
func checkIDs(idSets []IDSet) error {
	// Extract all unique IDVal
	idVals := IDSets(idSets).IDVals()

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

// WithMaxWorkers sets the maximum number of workers for rule evaluation
//   - -1:  sequential
//   - 0:   unlimited
//   - n>0: n workers
func (eng *Engine) WithMaxWorkers(n int) *Engine {
	eng.maxWorkers = n
	return eng
}

// WithMaxWorkers sets the maximum number of workers for rule evaluation
func (eng *Engine) getWorkerLimit(defaultLimit int) int {
	switch {
	case eng.maxWorkers <= -1:
		return 1 // sequential
	case eng.maxWorkers == 0:
		return defaultLimit // unlimited
	default:
		return eng.maxWorkers
	}
}
