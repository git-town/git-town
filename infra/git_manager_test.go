package infra

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestGitManagerCreateMemoizedEnvironment(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Errorf("cannot find temp dir: %s", err)
	}
	gm := NewGitManager(dir)
	err = gm.CreateMemoizedEnvironment()

	// verify error
	if err != nil {
		t.Errorf("creating memoized environment failed: %s", err)
	}

	// verify memoized folder exists
	memoizedPath := path.Join(dir, "memoized")
	if _, err := os.Stat(memoizedPath); os.IsNotExist(err) {
		t.Errorf("memoized directory (%s) not found", memoizedPath)
	}

}
