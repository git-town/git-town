package configdomain

import (
	"fmt"
)

// AliasableCommand defines Git Town commands that can shortened via Git aliases.
type AliasableCommand string

func (self AliasableCommand) Key() Key {
	switch self {
	case AliasableCommandAppend:
		return KeyAliasAppend
	case AliasableCommandDiffParent:
		return KeyAliasDiffParent
	case AliasableCommandHack:
		return KeyAliasHack
	case AliasableCommandKill:
		return KeyAliasKill
	case AliasableCommandPrepend:
		return KeyAliasPrepend
	case AliasableCommandPropose:
		return KeyAliasPropose
	case AliasableCommandRenameBranch:
		return KeyAliasRenameBranch
	case AliasableCommandRepo:
		return KeyAliasRepo
	case AliasableCommandShip:
		return KeyAliasShip
	case AliasableCommandSync:
		return KeyAliasSync
	}
	panic(fmt.Sprintf("don't know how to convert alias type %q into a config key", self))
}

func (self AliasableCommand) String() string { return string(self) }

const (
	AliasableCommandAppend       = AliasableCommand("append")
	AliasableCommandDiffParent   = AliasableCommand("diff-parent")
	AliasableCommandHack         = AliasableCommand("hack")
	AliasableCommandKill         = AliasableCommand("kill")
	AliasableCommandPrepend      = AliasableCommand("prepend")
	AliasableCommandPropose      = AliasableCommand("propose")
	AliasableCommandRenameBranch = AliasableCommand("rename-branch")
	AliasableCommandRepo         = AliasableCommand("repo")
	AliasableCommandShip         = AliasableCommand("ship")
	AliasableCommandSync         = AliasableCommand("sync")
)

// AliasableCommands provides all AliasType values.
func AliasableCommands() []AliasableCommand {
	return []AliasableCommand{
		AliasableCommandAppend,
		AliasableCommandDiffParent,
		AliasableCommandHack,
		AliasableCommandKill,
		AliasableCommandPrepend,
		AliasableCommandPropose,
		AliasableCommandRenameBranch,
		AliasableCommandRepo,
		AliasableCommandShip,
		AliasableCommandSync,
	}
}

func LookupAliasableCommand(key Key) *AliasableCommand {
	keyName := key.name
	for _, aliasableCommand := range AliasableCommands() {
		if aliasableCommand.String() == keyName {
			return &aliasableCommand
		}
	}
	return nil
}
