package fuzzy

import (
	"fugologic/crisp"
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

	// DefuzzificationSmallestOfMaxs returns the smallest of maximums
	DefuzzificationSmallestOfMaxs Defuzzification = func(fs Set, u crisp.Set) float64 {
		xSmallestMax, _ := defuzzificationMaximums(fs, u)
		return xSmallestMax
	}

	// DefuzzificationMiddleOfMaxs returns the middle of maximums
	DefuzzificationMiddleOfMaxs Defuzzification = func(fs Set, u crisp.Set) float64 {
		xSmallestMax, xLargestMax := defuzzificationMaximums(fs, u)
		return (xSmallestMax + xLargestMax) / 2
	}

	// DefuzzificationLargestOfMaxs returns the largest of maximums
	DefuzzificationLargestOfMaxs Defuzzification = func(fs Set, u crisp.Set) float64 {
		_, xLargestMax := defuzzificationMaximums(fs, u)
		return xLargestMax
	}
)

// defuzzificationMaximums returns the smallest of maximums and the largest of maximums
// E.g:
//
//	x = [0 1 2 3 4 5 6 7 8 9]
//	y = [0 0 1 1 2 2 1 1 0 0]
//	smallest of max is the left max for x=4 (y=2)
//	largest of max is the right max for x=5 (y=2)
func defuzzificationMaximums(fs Set, u crisp.Set) (float64, float64) {
	var xSmallestMax, xLargestMax float64
	var ySmallestMax, yLargestMax float64

	// Compute all y values
	values := u.Values()
	yValues := make([]float64, len(values))
	for i, x := range values {
		yValues[i] = fs(x)
	}

	// Find largest yi max for xi values where i in [0 ; n]
	// Find smallest yi max for xi values where i in [n ; 0]
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

// Aggregation represents the aggregation of 2 fuzzy set (for merging all result sets)
type Aggregation func(float64, float64) float64

var (
	AggregationUnion        Aggregation = union
	AggregationIntersection Aggregation = intersection
)

// defuzzer is responsible for collecting rule's results and to defuzz
type defuzzer struct {
	agg     Aggregation     // Aggregation of result fuzzy sets
	fct     Defuzzification // Defuzzification method
	results []IDSet         // From Val ID to list of result Set
}

// newDefuzzer builds a new Defuzzer instance
func newDefuzzer(fct Defuzzification, agg Aggregation, results []IDSet) defuzzer {
	return defuzzer{
		fct:     fct,
		agg:     agg,
		results: results,
	}
}

// defuzz the values
func (dfz defuzzer) defuzz() (DataOutput, error) {
	// Group IDSet by IDVal parent
	groups := make(map[*IDVal][]IDSet)
	universes := make(map[*IDVal]crisp.Set)
	for _, idSet := range dfz.results {
		idVal := idSet.parent
		groups[idVal] = append(groups[idVal], idSet)
		universes[idVal] = idVal.u
	}

	// For each group, apply defuzz
	values := make(DataOutput, len(dfz.results))
	for idVal, group := range groups {
		aggregation := dfz.aggregate(group)
		values[idVal] = dfz.fct(aggregation, universes[idVal])
	}
	return values, nil
}

// aggregate all sets into one (helper function): s = s1 U s2 U .. U sN
func (dfz defuzzer) aggregate(iss []IDSet) Set {
	result := iss[0].set
	for _, val := range iss[1:] {
		result = result.aggregate(val.set, dfz.agg)
	}
	return result
}
