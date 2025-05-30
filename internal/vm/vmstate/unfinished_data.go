package vmstate

import (
	"fmt"
	"strings"
	"time"

	"github.com/git-town/git-town/v21/internal/git/gitdomain"
)

// UnfinishedData has details about an unfinished run state.
type UnfinishedData struct {
	CanSkip   bool
	EndBranch gitdomain.LocalBranchName
	EndTime   time.Time
}

func (self UnfinishedData) String() string {
	result := strings.Builder{}
	result.WriteString("UnfinishedRunStateDetails {\n")
	result.WriteString("  CanSkip: ")
	result.WriteString(fmt.Sprintf("%t\n", self.CanSkip))
	result.WriteString("  EndBranch: ")
	result.WriteString(self.EndBranch.String())
	result.WriteRune('\n')
	return result.String()
}
