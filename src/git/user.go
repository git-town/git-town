package git

import "github.com/Originate/git-town/src/command"

// GetLocalAuthor returns the locally Git configured user
func GetLocalAuthor() string {
	name := command.Run("git", "config", "user.name").Output()
	email := command.Run("git", "config", "user.email").Output()
	return name + " <" + email + ">"
}
