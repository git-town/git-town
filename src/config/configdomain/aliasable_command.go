package configdomain

// AliasableCommand defines Git Town commands that can shortened via Git aliases.
type AliasableCommand string

func (self AliasableCommand) String() string { return string(self) }

const (
	AliasableCommandAppend       = AliasableCommand("append")
	AliasableCommandCompress     = AliasableCommand("compress")
	AliasableCommandContribute   = AliasableCommand("contribute")
	AliasableCommandDiffParent   = AliasableCommand("diff-parent")
	AliasableCommandHack         = AliasableCommand("hack")
	AliasableCommandKill         = AliasableCommand("kill")
	AliasableCommandObserve      = AliasableCommand("observe")
	AliasableCommandPark         = AliasableCommand("park")
	AliasableCommandPrepend      = AliasableCommand("prepend")
	AliasableCommandPropose      = AliasableCommand("propose")
	AliasableCommandRenameBranch = AliasableCommand("rename-branch")
	AliasableCommandRepo         = AliasableCommand("repo")
	AliasableCommandSetParent    = AliasableCommand("set-parent")
	AliasableCommandShip         = AliasableCommand("ship")
	AliasableCommandSync         = AliasableCommand("sync")
)
