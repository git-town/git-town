package configdomain

// AliasableCommand defines Git Town commands that can shortened via Git aliases.
type AliasableCommand string

// Key provides the key that configures this aliasable command in the Git config.
func (self AliasableCommand) Key() AliasKey {
	return AliasKey(AliasKeyPrefix + self.String())
}

func (self AliasableCommand) String() string { return string(self) }

const (
	// keep-sorted start
	AliasableCommandAppend     = AliasableCommand("append")
	AliasableCommandCompress   = AliasableCommand("compress")
	AliasableCommandContinue   = AliasableCommand("continue")
	AliasableCommandContribute = AliasableCommand("contribute")
	AliasableCommandDelete     = AliasableCommand("delete")
	AliasableCommandDiffParent = AliasableCommand("diff-parent")
	AliasableCommandDown       = AliasableCommand("down")
	AliasableCommandHack       = AliasableCommand("hack")
	AliasableCommandObserve    = AliasableCommand("observe")
	AliasableCommandPark       = AliasableCommand("park")
	AliasableCommandPrepend    = AliasableCommand("prepend")
	AliasableCommandPropose    = AliasableCommand("propose")
	AliasableCommandRename     = AliasableCommand("rename")
	AliasableCommandRepo       = AliasableCommand("repo")
	AliasableCommandSetParent  = AliasableCommand("set-parent")
	AliasableCommandShip       = AliasableCommand("ship")
	AliasableCommandSync       = AliasableCommand("sync")
	AliasableCommandUp         = AliasableCommand("up")
	// keep-sorted end
)
