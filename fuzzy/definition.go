package fuzzy

import (
	"errors"
	"fugologic/crisp"
	"fugologic/id"
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

// Not returns the complement of the current IDSet
func (is IDSet) Not() IDSet {
	return IDSet{
		set:    is.set.Complement(),
		uuid:   is.uuid,
		parent: is.parent,
	}
}

// Evaluate fetches the right input and returns the Set value
func (is IDSet) Evaluate(input DataInput) (float64, error) {
	x, err := input.find(is)
	if err != nil {
		return 0, err
	}
	return is.set(x), nil
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
		idSets: make(map[id.ID]IDSet),
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

// Get fuzzy set without checking content
func (iv IDVal) Get(name id.ID) IDSet {
	return iv.idSets[name]
}

// Fetch fuzzy set and check id presence
func (iv IDVal) Fetch(name id.ID) (IDSet, bool) {
	idSet, ok := iv.idSets[name]
	return idSet, ok
}
