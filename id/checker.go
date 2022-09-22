package id

import (
	"fmt"
)

// Identifiable represents an object with an ID
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
