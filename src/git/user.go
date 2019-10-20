package git

import "github.com/Originate/git-town/src/command"

// GetLocalAuthor returns the locally Git configured user
func GetLocalAuthor() string {
	name := command.Run("git", "config", "user.name").OutputSanitized()
	email := command.Run("git", "config", "user.email").OutputSanitized()
	return name + " <" + email + ">"
}
