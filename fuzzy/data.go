package fuzzy

import (
	"fmt"
)

// dataIO is a generic type to manipulate input/output values
// It represents couples of IDVal identifier and crisp value
type dataIO map[*IDVal]float64

// merge both data results and return a new filled structure
func (dt dataIO) merge(dt2 dataIO) dataIO {
	result := dataIO{}
	m := func(data dataIO) {
		for k, v := range data {
			result[k] = v
		}
	}

	m(dt)
	m(dt2)
	return result
}

// DataOutput represents the system output
type DataOutput dataIO

func (dout DataOutput) merge(dout2 DataOutput) DataOutput {
	result := dataIO(dout).merge(dataIO(dout2))
	return DataOutput(result)
}

// DataInput gives access to the input data system
type DataInput dataIO

func (din DataInput) merge(dout DataOutput) DataInput {
	result := dataIO(din).merge(dataIO(dout))
	return DataInput(result)
}

// value the linked IDVal to the IDSet, and then, its value
func (din DataInput) value(idSet IDSet) (float64, error) {
	if idSet.parent == nil {
		return 0, fmt.Errorf("input: cannot find parent for id set `%s`", idSet.uuid)
	}
	value, ok := din[idSet.parent]
	if !ok {
		return 0, fmt.Errorf("input: cannot find data for id val `%s` (id set `%s`)", idSet.parent.uuid, idSet.uuid)
	}
	return value, nil
}
