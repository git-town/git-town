//go:build unix

package git_test

import (
	"os"
	"syscall"
	"testing"

	"github.com/git-town/git-town/v24/pkg/asserts"
	"github.com/shoenig/test/must"
)

func openChangesTestFileOwnership(t *testing.T, path string) (int64, int64) {
	t.Helper()
	info := asserts.NoError1(os.Lstat(path))
	stat, isStat := info.Sys().(*syscall.Stat_t)
	must.True(t, isStat)
	return int64(stat.Uid), int64(stat.Gid)
}
