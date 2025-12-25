package runlog_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/state/runlog"
	"github.com/git-town/git-town/v22/pkg/asserts"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestEntry(t *testing.T) {
	t.Parallel()

	t.Run("serialize", func(t *testing.T) {
		t.Parallel()
		entry := runlog.Entry{
			Branches: map[gitdomain.BranchName]gitdomain.SHA{
				"main":            "111111",
				"origin/main":     "111111",
				"branch-1":        "222222",
				"origin/branch-1": "222222",
			},
			Command:        "git town sync",
			Event:          runlog.EventStart,
			Time:           time.Date(2025, 0o5, 28, 20, 34, 58, 123456789, time.UTC),
			PendingCommand: Some("sync"),
		}
		have := asserts.NoError1(json.MarshalIndent(entry, "", "  "))
		want := `
{
  "Branches": {
    "branch-1": "222222",
    "main": "111111",
    "origin/branch-1": "222222",
    "origin/main": "111111"
  },
  "Command": "git town sync",
  "Event": "start",
  "PendingCommand": "sync",
  "Time": "2025-05-28T20:34:58.123456789Z"
}`[1:]
		must.EqOp(t, want, string(have))
	})
}
