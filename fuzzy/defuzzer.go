package fuzzy

import (
	"fugologic/crisp"
	"math"
)

// Defuzzification method definition
// From a fuzzy set and its crisp values, evaluate only one crisp value
type Defuzzification func(fs Set, u crisp.Set) float64

var (
	// DefuzzificationCentroid is Sum(µ(xi)*xi) / Sum(µ(xi))
	DefuzzificationCentroid Defuzzification = defuzzificationCentroid

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

	// DefuzzificationBisector calculates the position under the curve where the areas on both sides are equal
	DefuzzificationBisector Defuzzification = defuzzificationBisector
)

func defuzzificationCentroid(fs Set, u crisp.Set) float64 {
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

func defuzzificationBisector(fs Set, u crisp.Set) float64 {
	values := u.Values()
	var left, right float64 // areas
	i := 0                  // left index (and result)
	j := len(values) - 1    // right index
	for ; i < j; i++ {      // move forward from start
		left += fs(values[i])               // increase left area
		for ; right <= left && i < j; j-- { // move backward from end
			right += fs(values[j]) // increase right area
		}
	}
	return values[i]
}

// Aggregation represents the aggregation of 2 fuzzy set (for merging all result sets)
type Aggregation func(float64, float64) float64

var (
	AggregationUnion        Aggregation = math.Max
	AggregationIntersection Aggregation = math.Min
)

// defuzzer is responsible for collecting rule's results and to defuzz
type defuzzer struct {
	agg Aggregation     // Aggregation of result fuzzy sets
	fct Defuzzification // Defuzzification method
}

// newDefuzzer builds a new Defuzzer instance
func newDefuzzer(fct Defuzzification, agg Aggregation) defuzzer {
	return defuzzer{
		fct: fct,
		agg: agg,
	}
}

// defuzz the values
func (dfz defuzzer) defuzz(iss []IDSet) DataOutput {
	// Group IDSet by IDVal parent
	groups := make(map[*IDVal][]IDSet)
	for _, idSet := range iss {
		idVal := idSet.parent
		groups[idVal] = append(groups[idVal], idSet)
	}

	// For each group, apply defuzz
	values := make(DataOutput, len(iss))
	for idVal, group := range groups {
		aggregation := dfz.aggregate(group)
		values[idVal] = dfz.fct(aggregation, idVal.u)
	}
	return values
}

// aggregate all sets into one (helper function): s = s1 U s2 U .. U sN
func (dfz defuzzer) aggregate(iss []IDSet) Set {
	result := iss[0].set
	for _, idSet := range iss[1:] {
		result = result.aggregate(idSet.set, dfz.agg)
	}
	return result
}
