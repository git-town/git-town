package git

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// Helpers

func getCurrentBranchNameDuringRebase() string {
	rawContent, err := ioutil.ReadFile(fmt.Sprintf("%s/.git/rebase-apply/head-name", GetRootDirectory()))
	if err != nil {
		// Git 2.26 introduces a new rebase backend, see https://github.com/git/git/blob/master/Documentation/RelNotes/2.26.0.txt
		rawContent, err = ioutil.ReadFile(fmt.Sprintf("%s/.git/rebase-merge/head-name", GetRootDirectory()))
		if err != nil {
			panic(err)
		}
	}
	content := strings.TrimSpace(string(rawContent))
	return strings.Replace(content, "refs/heads/", "", -1)
}
