package git

import "github.com/Originate/git-town/lib/util"

// GetLocalAuthor returns the locally Git configured user
func GetLocalAuthor() string {
	name := util.GetCommandOutput("git", "config", "user.name")
	email := util.GetCommandOutput("git", "config", "user.email")
	return name + " <" + email + ">"
}
