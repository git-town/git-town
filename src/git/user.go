package git

import "github.com/Originate/git-town/src/runner"

// GetLocalAuthor returns the locally Git configured user
func GetLocalAuthor() string {
	name := runner.New("git", "config", "user.name").Output()
	email := runner.New("git", "config", "user.email").Output()
	return name + " <" + email + ">"
}
