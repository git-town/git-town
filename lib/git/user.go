package git

import "github.com/Originate/git-town/lib/util"

func GetLocalAuthor() string {
	name := util.GetCommandOutput("git", "config", "user.name")
	email := util.GetCommandOutput("git", "config", "user.email")
	return name + " <" + email + ">"
}
