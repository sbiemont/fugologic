package id

import (
	"errors"
	"fmt"
)

// checker controls the Val and Set definition consistency
type Identifiable interface {
	ID() ID
}

// checker controls the Val and Set definition consistency
type checker []Identifiable

// NewChecker builds a new checker instance
func NewChecker(ids []Identifiable) checker {
	return ids
}

// Check launches the check operation on idSets
func (chk checker) Check() error {
	errPresence := chk.presence()
	if errPresence != nil {
		return errPresence
	}

	return chk.unicity()
}

// presence controls that all identifiers have been set
func (chk checker) presence() error {
	for _, identiable := range chk {
		if identiable.ID() == "" {
			return errors.New("id required")
		}
	}

	return nil
}

// unicity controls that all identifiers are unique
func (chk checker) unicity() error {
	uniqueIDs := make(map[ID]interface{})
	for _, identiable := range chk {
		uuid := identiable.ID()
		if _, exists := uniqueIDs[uuid]; exists {
			return fmt.Errorf("id `%s` already defined", uuid)
		}
		uniqueIDs[uuid] = nil
	}

	return nil
}
