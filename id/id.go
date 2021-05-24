package id

import "github.com/google/uuid"

// ID represents an identifier
type ID string

// NewID returns a random unique identifier
func NewID() ID {
	return ID(uuid.New().String())
}
