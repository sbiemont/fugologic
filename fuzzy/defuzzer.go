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

	// DefuzzificationSmallestOfMax returns the smallest of maximums
	DefuzzificationSmallestOfMax Defuzzification = func(fs Set, u crisp.Set) float64 {
		xSmallestMax, _ := defuzzificationMaximums(fs, u)
		return xSmallestMax
	}

	// DefuzzificationMiddleOfMax returns the middle of maximums
	DefuzzificationMiddleOfMax Defuzzification = func(fs Set, u crisp.Set) float64 {
		xSmallestMax, xLargestMax := defuzzificationMaximums(fs, u)
		return (xSmallestMax + xLargestMax) / 2
	}

	// DefuzzificationLargestOfMax returns the largest of maximums
	DefuzzificationLargestOfMax Defuzzification = func(fs Set, u crisp.Set) float64 {
		_, xLargestMax := defuzzificationMaximums(fs, u)
		return xLargestMax
	}
)

// defuzzificationMaximums returns the smallest of maximums and the largest of maximums
// E.g:
//  x = [0 1 2 3 4 5 6 7 8 9]
//  y = [0 0 1 1 2 2 1 1 0 0]
//  smallest of max is the left max for x=4 (y=2)
//  largest of max is the right max for x=5 (y=2)
func defuzzificationMaximums(fs Set, u crisp.Set) (float64, float64) {
	var xSmallestMax, xLargestMax float64
	var ySmallestMax, yLargestMax float64
	values := u.Values()
	yValues := make([]float64, len(values))
	for i, x := range values {
		yValues[i] = fs(x)
	}

	l := len(values) - 1
	for i, x := range values {
		x2 := values[l-i]
		y2 := yValues[l-i]
		y := yValues[i]
		if y >= yLargestMax {
			xLargestMax = x
			yLargestMax = y
		}
		if y2 >= ySmallestMax {
			xSmallestMax = x2
			ySmallestMax = y2
		}
	}
	return xSmallestMax, xLargestMax
}

// defuzzer is responsible for collecting rule's results and to defuzz
type defuzzer struct {
	fct     Defuzzification // Defuzzification method
	results []IDSet         // From Val ID to list of result Set
}

// newDefuzzer builds a new Defuzzer instance
func newDefuzzer(fct Defuzzification) defuzzer {
	return defuzzer{
		fct: fct,
	}
}

// add the sets to the defuzzer result
func (dfz *defuzzer) add(idSets []IDSet) {
	dfz.results = append(dfz.results, idSets...)
}

// defuzz the values
func (dfz defuzzer) defuzz() (DataOutput, error) {
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
func (defuzzer) union(iss []IDSet) Set {
	result := iss[0].set
	for _, val := range iss[1:] {
		result = result.Union(val.set)
	}
	return result
}
