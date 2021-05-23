package fuzzy

import (
	"errors"
	"fmt"

	"fugologic.git/id"
)

// checker controls the Val and Set definition consistency
type checker []IDSet

// newChecker builds a new Checker instance
func newChecker(idSets []IDSet) checker {
	return checker(idSets)
}

// check launches the check operation on idSets
func (chk checker) check() error {
	idVals, errPresence := chk.presence()
	if errPresence != nil {
		return errPresence
	}

	return chk.unicity(idVals)
}

// presence controls that all identifiers have been set
func (chk checker) presence() ([]*IDVal, error) {
	// Helper: check the ID presence
	identifiable := func(uuid id.ID) error {
		if uuid == "" {
			return errors.New("id required")
		}

		return nil
	}

	// Re-build the tree 1 Val => n Sets
	// Check presence
	idVals := make(map[*IDVal]interface{})
	for _, idSet := range chk {
		// Check parent presence
		parent := idSet.parent
		if parent == nil {
			return nil, fmt.Errorf("sets: no parent found for id set `%s`", idSet.uuid)
		}

		// Check Val id existence
		if err := identifiable(parent.uuid); err != nil {
			return nil, fmt.Errorf("values: %s", err)
		}

		// Check Set id existence
		if err := identifiable(idSet.uuid); err != nil {
			return nil, fmt.Errorf("sets: %s (for val id `%s`)", err, parent.uuid)
		}

		idVals[parent] = true
	}

	// To slice of unique IDVal
	var result []*IDVal
	for idVal := range idVals {
		result = append(result, idVal)
	}
	return result, nil
}

// unicity controls that all identifiers are unique
func (checker) unicity(idVals []*IDVal) error {
	// Helper: checks the ID unicity
	unique := func(mp map[id.ID]interface{}, uuid id.ID) error {
		if _, exists := mp[uuid]; exists {
			return fmt.Errorf("id `%s` already present", uuid)
		}
		mp[uuid] = true
		return nil
	}

	// Check unicity
	uniqueVals := make(map[id.ID]interface{})
	uniqueSets := make(map[id.ID]interface{})
	for _, idVal := range idVals {
		// Check unique Val id
		if err := unique(uniqueVals, idVal.uuid); err != nil {
			return fmt.Errorf("values: %s", err)
		}

		// Check unique Set id
		for _, idSet := range idVal.idSets {
			if err := unique(uniqueSets, idSet.uuid); err != nil {
				return fmt.Errorf("sets: %s (for val id `%s`)", err, idVal.uuid)
			}
		}
	}

	return nil
}
