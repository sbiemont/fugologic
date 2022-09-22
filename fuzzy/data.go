package fuzzy

import (
	"fmt"
)

// Data is a generic type to manipulate input/output values
// It represents couples of IDVal identifier and crisp value
type Data map[*IDVal]float64

// DataOutput represents the system output
type DataOutput Data

// DataInput gives access to the input data system
type DataInput Data

// find the linked IDVal to the IDSet, and then, its value
func (din DataInput) find(idSet IDSet) (float64, error) {
	if idSet.parent == nil {
		return 0, fmt.Errorf("input: cannot find parent for id set `%s`", idSet.uuid)
	}
	value, ok := din[idSet.parent]
	if !ok {
		return 0, fmt.Errorf("input: cannot find data for id val `%s` (id set `%s`)", idSet.parent.uuid, idSet.uuid)
	}
	return value, nil
}

// merge both data values and return the copy
func mergeData(d1, d2 map[*IDVal]float64) map[*IDVal]float64 {
	result := map[*IDVal]float64{}
	m := func(data map[*IDVal]float64) {
		for k, v := range data {
			result[k] = v
		}
	}

	m(d1)
	m(d2)
	return result
}
