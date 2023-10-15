package domain

import (
	"encoding/json"
	"fmt"
)

// SHA represents a Git SHA as a dedicated data type.
// This helps avoid stringly-typed code.
type SHA struct {
	id string
}

func EmptySHA() SHA {
	return SHA{
		id: "",
	}
}

// NewSHA creates a new SHA instance with the given value.
// The value is verified for correctness.
func NewSHA(id string) SHA {
	if !validateSHA(id) {
		panic(fmt.Sprintf("%q is not a valid Git SHA", id))
	}
	return SHA{id}
}

// validateSHA indicates whether the given SHA content is a valid Git SHA.
func validateSHA(content string) bool {
	if len(content) < 6 {
		return false
	}
	for _, c := range content {
		if c >= '0' && c <= '9' {
			continue
		}
		if c >= 'a' && c <= 'f' {
			continue
		}
		return false
	}
	return true
}

func (self SHA) IsEmpty() bool {
	return self.id == ""
}

// Location widens the type of this SHA to a more generic Location.
func (self SHA) Location() Location {
	return Location(self)
}

// MarshalJSON is used when serializing this SHA to JSON.
func (self SHA) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.id)
}

// Implementation of the fmt.Stringer interface.
func (self SHA) String() string { return self.id }

// TruncateTo provides a new SHA instance that contains a shorter checksum.
func (self SHA) TruncateTo(newLength int) SHA {
	if len(self.id) < newLength {
		return self
	}
	return NewSHA(self.id[0:newLength])
}

// UnmarshalJSON is used when de-serializing JSON into a SHA.
func (self *SHA) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &self.id)
}
