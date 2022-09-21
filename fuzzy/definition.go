package fuzzy

import (
	"fmt"

	"fugologic/crisp"
	"fugologic/id"
)

// IDSet represents a static Set with an ID
type IDSet struct {
	set    Set    // membership function
	uuid   id.ID  // identifier is only used to have error information
	parent *IDVal // parent leads to the IDVal (for defuzzification)
}

// NewIDSet associates a Set with an UUID
func NewIDSet(set Set, parent *IDVal) IDSet {
	return NewIDSetCustom(id.NewID(), set, parent)
}

// NewIDSetCustom associates a Set with a custom ID
func NewIDSetCustom(uuid id.ID, set Set, parent *IDVal) IDSet {
	idSet := IDSet{
		set:    set,
		uuid:   uuid,
		parent: parent,
	}
	parent.add(&idSet)
	return idSet
}

// ID returns the identifier
func (is IDSet) ID() id.ID {
	return is.uuid
}

// Parent returns the fuzzy value of the current fuzzy set
func (is IDSet) Parent() *IDVal {
	return is.parent
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
	idSets []*IDSet // only used for consistency checking
}

// NewIDVal associates a list of Set with an UUID
func NewIDVal(u crisp.Set) IDVal {
	return NewIDValCustom(id.NewID(), u)
}

// NewIDValCustom associates a list of Set with a custom ID
func NewIDValCustom(uuid id.ID, u crisp.Set) IDVal {
	return IDVal{
		uuid: uuid,
		u:    u,
	}
}

// add a new IDSet (only used for control)
func (iv *IDVal) add(idSet *IDSet) {
	iv.idSets = append(iv.idSets, idSet)
}

// ID returns the identifier
func (iv IDVal) ID() id.ID {
	return iv.uuid
}

// U returns the crisp set of values
func (iv IDVal) U() crisp.Set {
	return iv.u
}

// checkIDs of a list of IDSet
// Get all unique IDVal, check them and their whole IDSet
func checkIDs(idSets []IDSet) error {
	// Extract all unique IDVal
	var idVals = make(map[*IDVal]interface{})
	for _, idSet := range idSets {
		idVals[idSet.parent] = nil
	}

	// Extract all IDSets of IDVals and compare them
	var vals []id.Identifiable
	var sets []id.Identifiable
	for idVal := range idVals {
		vals = append(vals, idVal)
		for _, idSet := range idVal.idSets {
			sets = append(sets, idSet)
		}
	}

	// Check all
	if err := id.NewChecker(vals).Check(); err != nil {
		return fmt.Errorf("values: %s", err)
	}
	if err := id.NewChecker(sets).Check(); err != nil {
		return fmt.Errorf("sets: %s", err)
	}
	return nil
}
