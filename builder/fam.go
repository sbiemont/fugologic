package builder

import (
	"fmt"
	"fugologic/fuzzy"
	"fugologic/id"
)

// FuzzyAssoMatrix (or FAM) stands for "Fuzzy Associative Matrix"
// It is a compact way to express fuzzy logic rules in tabular form
// https://en.wikipedia.org/wiki/fuzzy_associative_matrix
type FuzzyAssoMatrix struct {
	cfg famConfig

	rules []fuzzy.Rule // stores internal rules (like the builder)
}

// famConfig gathers data for processing the rules
type famConfig struct {
	and    fuzzy.Connector
	impl   fuzzy.Implication
	agg    fuzzy.Aggregation
	defuzz fuzzy.Defuzzification
}

// NewFuzzyAssoMatrix create a new FAM instance
func NewFuzzyAssoMatrix(
	and fuzzy.Connector,
	impl fuzzy.Implication,
	agg fuzzy.Aggregation,
	defuzz fuzzy.Defuzzification,
) FuzzyAssoMatrix {
	return FuzzyAssoMatrix{
		cfg: famConfig{
			and:    and,
			impl:   impl,
			agg:    agg,
			defuzz: defuzz,
		},
	}
}

// add a new rule to the matrix
func (fam *FuzzyAssoMatrix) add(rule fuzzy.Rule) {
	fam.rules = append(fam.rules, rule)
}

// Engine builds an engine from the configuration
func (fam FuzzyAssoMatrix) Engine() (fuzzy.Engine, error) {
	return fuzzy.NewEngine(fam.rules, fam.cfg.agg, fam.cfg.defuzz)
}

// famExpression is an internal structure to store redundant information
type famExpression struct {
	val *fuzzy.IDVal
	ids []id.ID
}

// ifExpression is an internal structure to store builder + first part of the rule
type ifExpression struct {
	fam   *FuzzyAssoMatrix
	ifExp famExpression
}

// andExpression is an internal structure to store builder + first part of the rule + second part of the rule
type andExpression struct {
	fam    *FuzzyAssoMatrix
	ifExp  famExpression
	andExp famExpression
}

type famValues struct {
	fam     *FuzzyAssoMatrix
	ifVal   *fuzzy.IDVal
	andVal  *fuzzy.IDVal
	thenVal *fuzzy.IDVal
}

// Asso defines the rules pattern of fuzzy values
// if <a> and <b> then <c>
func (fam *FuzzyAssoMatrix) Asso(ifVal, andVal, thenVal *fuzzy.IDVal) famValues {
	return famValues{
		fam:     fam,
		ifVal:   ifVal,
		andVal:  andVal,
		thenVal: thenVal,
	}
}

// Matrix creates all rules by merging input <if> with input <and> into output <then>
// The algo is:
//
// - For i: values of a
//   - For j: values of b
//     Rule = if (a[i]) and (b[j]) then (c[i][j])
func (fv famValues) Matrix(ifSets []id.ID, andThenSets map[id.ID][]id.ID) error {
	// Fetch for id-set within an id-val
	// Use a map to fetch data only once
	type mapID map[id.ID]fuzzy.IDSet
	fetch := func(info string, v *fuzzy.IDVal, uuid id.ID, mp mapID) (fuzzy.IDSet, error) {
		// Fetch data from map
		premise, exists := mp[uuid]
		if exists {
			return premise, nil
		}
		// Fetch data from value
		premise, ok := v.Fetch(uuid)
		if !ok {
			return fuzzy.IDSet{}, fmt.Errorf("'%s' statement, cannot find %s from %s", info, uuid, v.ID())
		}
		// Store and return data
		mp[uuid] = premise
		return premise, nil
	}

	// Control sizes
	checkSize := func(actual, expected int) error {
		if actual != expected {
			return fmt.Errorf("rule, sizes should be the same (found: %d, expected: %d)", actual, expected)
		}
		return nil
	}

	ifMap := make(mapID)
	andMap := make(mapID)
	thenMap := make(mapID)

	n := len(ifSets)
	for i, ifID := range ifSets {
		ifSet, errIf := fetch("if", fv.ifVal, ifID, ifMap)
		if errIf != nil {
			return errIf
		}

		for andID, thenSets := range andThenSets {
			errSize := checkSize(len(thenSets), n)
			if errSize != nil {
				return errSize
			}
			andSet, errAnd := fetch("and", fv.andVal, andID, andMap)
			if errAnd != nil {
				return errAnd
			}
			if thenSets[i].Empty() {
				// No rule is the "then" statement is empty
				continue
			}
			thenSet, errThen := fetch("then", fv.thenVal, thenSets[i], thenMap)
			if errThen != nil {
				return errThen
			}

			fv.fam.add(fuzzy.NewRule(
				fuzzy.NewExpression([]fuzzy.Premise{ifSet, andSet}, fv.fam.cfg.and),
				fv.fam.cfg.impl,
				[]fuzzy.IDSet{thenSet},
			))
		}
	}

	if len(ifMap) != n {
		return fmt.Errorf("'if' statement, duplicated headers found")
	}

	return nil
}
