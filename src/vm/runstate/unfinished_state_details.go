package runstate

import (
	"fmt"
	"strings"
	"time"

	"github.com/git-town/git-town/v13/src/git/gitdomain"
)

// UnfinishedRunStateDetails has details about an unfinished run state.
type UnfinishedRunStateDetails struct {
	CanSkip   bool
	EndBranch gitdomain.LocalBranchName
	EndTime   time.Time
}

func (self UnfinishedRunStateDetails) String() string {
	result := strings.Builder{}
	result.WriteString("UnfinishedRunStateDetails {\n")
	result.WriteString("  CanSkip: ")
	result.WriteString(fmt.Sprintf("%t\n", self.CanSkip))
	result.WriteString("  EndBranch: ")
	result.WriteString(self.EndBranch.String())
	result.WriteRune('\n')
	return result.String()
}
