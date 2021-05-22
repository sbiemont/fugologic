package fuzzy

import (
	"fugologic.git/crisp"
	"fugologic.git/id"
)

// Defuzzification method definition
// From a fuzzy set and its crisp values, evaluate only one crisp value
type Defuzzification func(fs Set, u crisp.Set) float64

var (
	// DefuzzificationCentroid is Sum(µ(xi)*xi) / Sum(µ(xi))
	DefuzzificationCentroid Defuzzification = func(fs Set, u crisp.Set) float64 {
		var mx, m float64
		for _, x := range u.Values() {
			y := fs(x)
			mx += y * x
			m += y
		}

		if m == 0 {
			return 0
		}

		return mx / m
	}

	// TODO add more defuzzification method

)

// Defuzzer is responsible for collecting rule's results and to defuzz
type Defuzzer struct {
	fct     Defuzzification // Defuzzification method
	results []IDSet         // From Val ID to list of result Set
}

// NewDefuzzer builds a new Defuzzer instance
func NewDefuzzer(fct Defuzzification) Defuzzer {
	return Defuzzer{
		fct: fct,
	}
}

// Add merges the sets to the defuzzer result
func (dfz *Defuzzer) Add(idSets []IDSet) {
	dfz.results = append(dfz.results, idSets...)
}

// Defuzz the values
func (dfz Defuzzer) Defuzz() (DataOutput, error) {
	// Group IDSet by IDVal parent
	groups := make(map[id.ID][]IDSet)
	universes := make(map[id.ID]crisp.Set)
	for _, idSet := range dfz.results {
		idVal := idSet.parent
		groups[idVal.uuid] = append(groups[idVal.uuid], idSet)
		universes[idVal.uuid] = idVal.u
	}

	// For each group, apply defuzz
	values := make(map[id.ID]float64, len(dfz.results))
	for id, group := range groups {
		union := dfz.union(group)
		values[id] = dfz.fct(union, universes[id])
	}
	return values, nil
}

// union all sets into one (helper function): s = s1 U s2 U .. U sN
func (Defuzzer) union(iss []IDSet) Set {
	result := iss[0].set
	for _, val := range iss[1:] {
		result = result.Union(val.set)
	}
	return result
}
