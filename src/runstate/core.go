// Package runstate represents the current execution status of a Git Town command,
// i.e. which steps to execute via the `runvm`.
package runstate

import (
	"fmt"
	"strings"
	"time"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/gohacks"
	"github.com/git-town/git-town/v9/src/step"
)

// UnfinishedRunStateDetails has details about an unfinished run state.
type UnfinishedRunStateDetails struct {
	CanSkip   bool
	EndBranch domain.LocalBranchName
	EndTime   time.Time
}

func (ursd UnfinishedRunStateDetails) String() string {
	result := strings.Builder{}
	result.WriteString("UnfinishedRunStateDetails {\n")
	result.WriteString("  CanSkip: ")
	result.WriteString(fmt.Sprintf("%t\n", ursd.CanSkip))
	result.WriteString("  EndBranch: ")
	result.WriteString(ursd.EndBranch.String())
	result.WriteString("\n")
	return result.String()
}

func isCheckoutStep(step step.Step) bool {
	typeName := gohacks.TypeName(step)
	return typeName == "Checkout" || typeName == "CheckoutIfExists"
}
