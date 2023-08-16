package domain

import (
	"encoding/json"
	"fmt"
)

// SHA represents a Git SHA as a dedicated data type.
// This helps avoid stringly-typed code.
type SHA struct {
	Location
}

// NewSHA creates a new SHA instance with the given value.
// The value is verified for correctness.
func NewSHA(id string) SHA {
	if !validateSHA(id) {
		panic(fmt.Sprintf("%q is not a valid Git SHA", id))
	}
	return SHA{Location{id}}
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

func (s SHA) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.value)
}

// Implements the fmt.Stringer interface.
func (s SHA) String() string { return s.value }

// TruncateTo provides a new SHA instance that contains a shorter checksum.
func (s SHA) TruncateTo(newLength int) SHA {
	return NewSHA(s.value[0:newLength])
}

func (s *SHA) UnmarshalJSON(b []byte) error {
	var t string
	err := json.Unmarshal(b, &t)
	if err != nil {
		return err
	}
	s.value = t
	return nil
}
