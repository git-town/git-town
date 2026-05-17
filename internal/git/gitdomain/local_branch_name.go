package gitdomain

import (
	"github.com/git-town/git-town/v23/internal/gohacks/stringss"
	. "github.com/git-town/git-town/v23/pkg/prelude"
)

// LocalBranchName is the name of a local Git branch.
// The zero value is an empty local branch name,
// i.e. a local branch name that is unknown or not configured.
type LocalBranchName stringss.Trimmed

func LocalBranchNameOrPanic(id stringss.Trimmed) LocalBranchName {
	if !isValidLocalBranchName(id) {
		panic("local branch names cannot be empty")
	}
	return LocalBranchName(id)
}

func LocalBranchNameOpt(id stringss.Trimmed) Option[LocalBranchName] {
	if isValidLocalBranchName(id) {
		return Some(LocalBranchName(id))
	}
	return None[LocalBranchName]()
}

func isValidLocalBranchName(value stringss.Trimmed) bool {
	return len(value) > 0
}

// BranchName widens the type of this LocalBranchName to a more generic BranchName.
func (self LocalBranchName) BranchName() BranchName {
	return BranchName(string(self))
}

// Location widens the type of this LocalBranchName to a more generic Location.
func (self LocalBranchName) Location() Location {
	return NewLocation(string(self))
}

// RefName provides the fully qualified reference name for this branch.
func (self LocalBranchName) RefName() string {
	return "refs/heads/" + self.String()
}

// String implements the fmt.Stringer interface.
func (self LocalBranchName) String() string { return string(self) }
