package git

import "github.com/Originate/git-town/src/command"

// GetLocalAuthor returns the locally Git configured user
func GetLocalAuthor() string {
	name := command.MustRun("git", "config", "user.name").OutputSanitized()
	email := command.MustRun("git", "config", "user.email").OutputSanitized()
	return name + " <" + email + ">"
}
