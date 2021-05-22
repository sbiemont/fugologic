package fuzzy

import (
	"fmt"

	"fugologic.git/id"
)

// DataOutput represents the system output
type DataOutput map[id.ID]float64

// DataInput gives access to the input data system
type DataInput map[id.ID]float64

// find the linked IDVal to the IDSet, and then, its value
func (din DataInput) find(idSet IDSet) (float64, error) {
	if idSet.parent == nil {
		return 0, fmt.Errorf("input: cannot find parent for id set `%s`", idSet.uuid)
	}
	value, ok := din[idSet.parent.uuid]
	if !ok {
		return 0, fmt.Errorf("input: cannot find data for id val `%s` (id set `%s`)", idSet.parent.uuid, idSet.uuid)
	}
	return value, nil
}
