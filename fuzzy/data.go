package fuzzy

import (
	"fmt"
)

// Data is a generic type to manipulate input/output values
// It represents couples of IDVal identifier and crisp value
type Data map[*IDVal]float64

// merge both data results and return a new filled structure
func (dt Data) merge(dt2 Data) Data {
	result := Data{}
	m := func(data Data) {
		for k, v := range data {
			result[k] = v
		}
	}

	m(dt)
	m(dt2)
	return result
}

// DataOutput represents the system output
type DataOutput Data

func (dout DataOutput) merge(dout2 DataOutput) DataOutput {
	result := Data(dout).merge(Data(dout2))
	return DataOutput(result)
}

// DataInput gives access to the input data system
type DataInput Data

func (din DataInput) merge(dout DataOutput) DataInput {
	result := Data(din).merge(Data(dout))
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
