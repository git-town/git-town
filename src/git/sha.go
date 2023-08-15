package git

import (
	"encoding/json"
	"fmt"
)

// SHA represents a Git SHA as a dedicated data type.
// This helps avoid stringly-typed code.
type SHA struct {
	content string
}

// NewSHA creates a new SHA instance with the given value.
// The value is verified for correctness.
func NewSHA(content string) SHA {
	if !validateSHA(content) {
		panic(fmt.Sprintf("%q is not a valid Git SHA", content))
	}
	return SHA{content}
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

// ErrorSHA provides the zero value for a Git SHA, to be used only when returning a SHA that should be ignored because it is returned as part of an error.
// This is needed because we need to return invalid SHAs from functions that return an error.
func ErrorSHA() SHA {
	return SHA{""}
}

func (s SHA) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.content)
}

// Implements the fmt.Stringer interface for Git SHAs.
func (s SHA) String() string { return s.content }

// TruncateTo provides a new SHA instance that contains a shorter checksum.
func (s SHA) TruncateTo(newLength int) SHA {
	return SHA{s.content[0:newLength]}
}

func (s *SHA) UnmarshalJSON(b []byte) error {
	var t string
	err := json.Unmarshal(b, &t)
	if err != nil {
		return err
	}
	*s = SHA{t}
	return nil
}
