package gh

import "os/exec"

func Detect(querier Querier) bool {
	// detect gh executable
	ghPath, err := exec.LookPath("gh")
	if err != nil {
		return false
	}
}
