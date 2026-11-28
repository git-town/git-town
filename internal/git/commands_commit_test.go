package git_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/git-town/git-town/v24/internal/config/configdomain"
	"github.com/git-town/git-town/v24/internal/git/gitdomain"
	"github.com/git-town/git-town/v24/internal/gohacks/stringss"
	testsubshell "github.com/git-town/git-town/v24/internal/test/subshell"
	"github.com/git-town/git-town/v24/internal/test/testruntime"
	"github.com/git-town/git-town/v24/pkg/asserts"
	. "github.com/git-town/git-town/v24/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestCommitMessageVerbatimDecodesCommitEncoding(t *testing.T) {
	t.Parallel()
	runtime := testruntime.Create(t)
	runtime.CreateFile("file", "content")
	runtime.MustRun("git", "add", "file")
	tree := runtime.MustQueryTrim("git", "write-tree")
	runtime.MustRun("git", "config", "i18n.commitEncoding", "ISO-8859-1")
	output := runtime.MustQueryWith(&testsubshell.Options{Input: Some("ol\xe1 encoded"), TTY: true}, "git", "commit-tree", tree.String())
	sha := gitdomain.NewSHAOrPanic(stringss.Trim(output))
	have := asserts.NoError1(runtime.Git.CommitMessageVerbatim(runtime.TestRunner, sha))
	want := gitdomain.CommitMessage("olá encoded")
	must.EqOp(t, want, have)
}

func TestCommitMessageVerbatimIgnoresSignatureOutput(t *testing.T) {
	t.Parallel()
	runtime := testruntime.Create(t)
	runtime.CreateFile("file", "content")
	runtime.MustRun("git", "add", "file")
	tree := runtime.MustQueryTrim("git", "write-tree")
	commitContent := fmt.Sprintf("tree %s\nauthor Test <test@example.com> 1764300000 +0000\ncommitter Test <test@example.com> 1764300000 +0000\ngpgsig -----BEGIN PGP SIGNATURE-----\n \nabcdef\n -----END PGP SIGNATURE-----\n\nsigned message", tree)
	output := runtime.MustQueryWith(&testsubshell.Options{Input: Some(commitContent), TTY: true}, "git", "hash-object", "-t", "commit", "-w", "--stdin")
	sha := gitdomain.NewSHAOrPanic(stringss.Trim(output))
	runtime.MustRun("git", "config", "gpg.program", "git")
	runtime.MustRun("git", "config", "log.showSignature", "true")
	have := asserts.NoError1(runtime.Git.CommitMessageVerbatim(runtime.TestRunner, sha))
	want := gitdomain.CommitMessage("signed message")
	must.EqOp(t, want, have)
}

func TestCommitMessageVerbatimRejectsInvalidUTF8(t *testing.T) {
	t.Parallel()
	runtime := testruntime.Create(t)
	runtime.CreateFile("file", "content")
	runtime.MustRun("git", "add", "file")
	tree := runtime.MustQueryTrim("git", "write-tree")
	commitContent := fmt.Sprintf("tree %s\nauthor Test <test@example.com> 1764300000 +0000\ncommitter Test <test@example.com> 1764300000 +0000\nencoding ISO-2022-JP\n\ninvalid ISO-2022-JP \xff", tree)
	output := runtime.MustQueryWith(&testsubshell.Options{Input: Some(commitContent), TTY: true}, "git", "hash-object", "-t", "commit", "-w", "--stdin")
	sha := gitdomain.NewSHAOrPanic(stringss.Trim(output))
	_, err := runtime.Git.CommitMessageVerbatim(runtime.TestRunner, sha)
	must.ErrorContains(t, err, "UTF-8")
	must.ErrorContains(t, err, sha.String())
}

func TestCommitVerbatimLargeMessageWithoutFinalNewline(t *testing.T) {
	t.Parallel()
	runtime := testruntime.Create(t)
	runtime.CreateFile("file", "content")
	runtime.MustRun("git", "add", "file")
	want := gitdomain.CommitMessage(strings.Repeat("x", 256*1024))
	err := runtime.Git.Commit(runtime.TestRunner, configdomain.UseCustomMessageVerbatim(want), None[gitdomain.Author](), configdomain.CommitHookDisabled)
	must.NoError(t, err)
	sha := gitdomain.NewSHAOrPanic(runtime.MustQueryTrim("git", "rev-parse", "HEAD"))
	have := asserts.NoError1(runtime.Git.CommitMessageVerbatim(runtime.TestRunner, sha))
	must.Eq(t, len(want), len(have))
	must.True(t, want == have)
}

func TestCommitVerbatimPreservesWhitespaceAndMissingFinalNewline(t *testing.T) {
	t.Parallel()
	runtime := testruntime.Create(t)
	runtime.CreateFile("file", "content")
	runtime.MustRun("git", "add", "file")
	want := gitdomain.CommitMessage("commit 1\n\n# preserved comment\ntrailing space  ")
	err := runtime.Git.Commit(runtime.TestRunner, configdomain.UseCustomMessageVerbatim(want), None[gitdomain.Author](), configdomain.CommitHookDisabled)
	must.NoError(t, err)
	sha := gitdomain.NewSHAOrPanic(runtime.MustQueryTrim("git", "rev-parse", "HEAD"))
	have := asserts.NoError1(runtime.Git.CommitMessageVerbatim(runtime.TestRunner, sha))
	must.EqOp(t, want, have)
}

func TestCommitVerbatimRoundTripWritesUTF8WithLegacyCommitEncoding(t *testing.T) {
	t.Parallel()
	runtime := testruntime.Create(t)
	runtime.CreateFile("file", "content")
	runtime.MustRun("git", "add", "file")
	tree := runtime.MustQueryTrim("git", "write-tree")
	runtime.MustRun("git", "config", "i18n.commitEncoding", "ISO-8859-1")
	output1 := runtime.MustQueryWith(&testsubshell.Options{Input: Some("ol\xe1 encoded 1"), TTY: true}, "git", "commit-tree", tree.String())
	sha1 := gitdomain.NewSHAOrPanic(stringss.Trim(output1))
	output2 := runtime.MustQueryWith(&testsubshell.Options{Input: Some("ol\xe1 encoded 2"), TTY: true}, "git", "commit-tree", tree.String())
	sha2 := gitdomain.NewSHAOrPanic(stringss.Trim(output2))
	message1 := asserts.NoError1(runtime.Git.CommitMessageVerbatim(runtime.TestRunner, sha1))
	message2 := asserts.NoError1(runtime.Git.CommitMessageVerbatim(runtime.TestRunner, sha2))
	want := gitdomain.CommitMessages{message1, message2}.Join()
	err := runtime.Git.Commit(runtime.TestRunner, configdomain.UseCustomMessageVerbatim(want), None[gitdomain.Author](), configdomain.CommitHookDisabled)
	must.NoError(t, err)
	sha := gitdomain.NewSHAOrPanic(runtime.MustQueryTrim("git", "rev-parse", "HEAD"))
	have := asserts.NoError1(runtime.Git.CommitMessageVerbatim(runtime.TestRunner, sha))
	must.EqOp(t, want, have)
}
