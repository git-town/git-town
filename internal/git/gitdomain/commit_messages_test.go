package gitdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v24/internal/git/gitdomain"
	"github.com/shoenig/test/must"
)

func TestCommitMessagesJoinEmptyList(t *testing.T) {
	t.Parallel()
	give := gitdomain.CommitMessages{}
	have := give.Join()
	want := gitdomain.CommitMessage("")
	must.EqOp(t, want, have)
}

func TestCommitMessagesJoinMissingFinalNewline(t *testing.T) {
	t.Parallel()
	give := gitdomain.NewCommitMessages("commit 1", "commit 2")
	have := give.Join()
	want := gitdomain.CommitMessage("commit 1\n\ncommit 2")
	must.EqOp(t, want, have)
}

func TestCommitMessagesJoinPreservesWhitespace(t *testing.T) {
	t.Parallel()
	give := gitdomain.CommitMessages{"", "commit  "}
	have := give.Join()
	want := gitdomain.CommitMessage("\n\ncommit  ")
	must.EqOp(t, want, have)
}

func TestCommitMessagesJoinWithFinalNewlines(t *testing.T) {
	t.Parallel()
	give := gitdomain.NewCommitMessages("commit 1\n", "commit 2\n")
	have := give.Join()
	want := gitdomain.CommitMessage("commit 1\n\ncommit 2\n")
	must.EqOp(t, want, have)
}

func TestCommitMessagesJoinWithTrailingBlankLine(t *testing.T) {
	t.Parallel()
	give := gitdomain.NewCommitMessages("commit 1\n\n", "commit 2")
	have := give.Join()
	want := gitdomain.CommitMessage("commit 1\n\n\ncommit 2")
	must.EqOp(t, want, have)
}
