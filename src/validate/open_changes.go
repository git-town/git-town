package validate

import "fmt"

func NoOpenChanges() error {
	hasOpenChanges, err := pr.Backend.HasOpenChanges()
	if err != nil {
		return err
	}
	if hasOpenChanges {
		err = fmt.Errorf("you have uncommitted changes. Did you mean to commit them before shipping?")
		return
	}
}
