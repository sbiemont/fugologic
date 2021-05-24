package fuzzy

import (
	"fugologic.git/crisp"
	"fugologic.git/id"
)

// IDSet represents a static Set with an ID
type IDSet struct {
	set    Set
	uuid   id.ID
	parent *IDVal
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

// And returns the expression IDSet AND Premise
func (is IDSet) And(premise Premise) Premise {
	return NewExpression([]Premise{is, premise}, ConnectorZadehAnd)
}

// Or returns the expression IDSet OR Premise
func (is IDSet) Or(premise Premise) Premise {
	return NewExpression([]Premise{is, premise}, ConnectorZadehOr)
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
