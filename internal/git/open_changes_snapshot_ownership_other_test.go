//go:build !unix

package git_test

import "testing"

func openChangesTestFileOwnership(_ *testing.T, _ string) (int64, int64) {
	return 0, 0
}
