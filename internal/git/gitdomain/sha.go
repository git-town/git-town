package gitdomain

import (
	"fmt"
)

// SHA represents a Git SHA as a dedicated data type.
// This helps avoid stringly-typed code.
type SHA string

// NewSHA creates a new SHA instance with the given value.
// The value is verified for correctness.
func NewSHA(id string) SHA {
	result, err := NewSHAErr(id)
	if err != nil {
		panic(err)
	}
	return result
}

func NewSHAErr(id string) (SHA, error) {
	if !validateSHA(id) {
		return SHA(""), fmt.Errorf("%q is not a valid Git SHA", id)
	}
	return SHA(id), nil
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

// Location widens the type of this SHA to a more generic Location.
func (self SHA) Location() Location {
	return Location(string(self))
}

// Implementation of the fmt.Stringer interface.
func (self SHA) String() string {
	return string(self)
}

// Truncates reduces the length of this SHA to the given length.
// This is only for test code.
// Please use git.Commands.ShortenSHA in production code.
func (self SHA) Truncate(newLen int) SHA {
	return NewSHA(self.String()[:newLen])
}
