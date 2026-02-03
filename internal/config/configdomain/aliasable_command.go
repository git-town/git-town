package configdomain

// AliasableCommand defines Git Town commands that can shortened via Git aliases.
type AliasableCommand string

// Key provides the key that configures this aliasable command in the Git config.
func (self AliasableCommand) Key() AliasKey {
	return AliasKey(AliasKeyPrefix + self.String())
}

func (self AliasableCommand) String() string { return string(self) }

const (
	AliasableCommandAppend     = AliasableCommand("append")
	AliasableCommandCompress   = AliasableCommand("compress")
	AliasableCommandContinue   = AliasableCommand("continue")
	AliasableCommandContribute = AliasableCommand("contribute")
	AliasableCommandDiffParent = AliasableCommand("diff-parent")
	AliasableCommandDown       = AliasableCommand("down")
	AliasableCommandHack       = AliasableCommand("hack")
	AliasableCommandDelete     = AliasableCommand("delete")
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
)
