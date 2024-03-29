package fuzzy

import (
	"errors"

	"github.com/sbiemont/fugologic/crisp"
	"github.com/sbiemont/fugologic/id"
)

// IDSet represents a static Set with an ID
type IDSet struct {
	set    Set    // membership function
	uuid   id.ID  // identifier is only used to have error information
	parent *IDVal // parent leads to the IDVal (for defuzzification)
}

// ID returns the identifier
func (is IDSet) ID() id.ID {
	return is.uuid
}

// Evaluate fetches the right input and returns the Set value
func (is IDSet) Evaluate(input DataInput) (float64, error) {
	x, err := input.value(is)
	if err != nil {
		return 0, err
	}
	return is.set(x), nil
}

// IDSets is an helper for managing a list of IDSet
type IDSets []IDSet

// IDVals extract unique IDVal from the list of IDSet
func (iss IDSets) IDVals() map[*IDVal]struct{} {
	result := make(map[*IDVal]struct{})
	for _, idSet := range iss {
		result[idSet.parent] = struct{}{}
	}
	return result
}

// IDVal represents a static Val with an ID and a set of crisp values
type IDVal struct {
	uuid   id.ID
	u      crisp.Set
	idSets map[id.ID]IDSet // used for consistency checking
}

// NewIDVal associates a list of Set with a custom ID
func NewIDVal(
	uuid id.ID, // uuid is the identifier of the fuzzy value (empty uuid is rejected)
	u crisp.Set, // u is the crisp universe of the value
	sets map[id.ID]Set, // sets are the list of couples uuid + fuzzy set (empty uuids are rejected)
) (*IDVal, error) {
	// Check for empty id
	if uuid.Empty() {
		return nil, errors.New("id val cannot be empty")
	}
	iv := &IDVal{
		uuid:   uuid,
		u:      u,
		idSets: make(map[id.ID]IDSet, len(sets)),
	}

	// Convert to id set and check for empty id
	for name, set := range sets {
		if name.Empty() {
			return nil, errors.New("id set cannot be empty")
		}
		iv.idSets[name] = IDSet{
			set:    set,
			uuid:   name,
			parent: iv,
		}
	}
	return iv, nil
}

// ID returns the identifier
func (iv IDVal) ID() id.ID {
	return iv.uuid
}

// U retrieves the crisp data universe
func (iv IDVal) U() crisp.Set {
	return iv.u
}

// Get fuzzy set without checking content
func (iv IDVal) Get(name id.ID) IDSet {
	return iv.idSets[name]
}

// Fetch fuzzy set and check id presence
func (iv IDVal) Fetch(name id.ID) (IDSet, bool) {
	idSet, ok := iv.idSets[name]
	return idSet, ok
}
