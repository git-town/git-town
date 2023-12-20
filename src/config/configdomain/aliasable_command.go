package configdomain

import (
	"fmt"
	"strings"
)

// AliasableCommand defines Git Town commands that can shortened via Git aliases.
// This is a type-safe enum, see https://npf.io/2022/05/safer-enums.
type AliasableCommand struct {
	name string
}

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

func (self AliasableCommand) String() string { return self.name }

var (
	AliasableCommandAppend       = AliasableCommand{"append"}        //nolint:gochecknoglobals
	AliasableCommandDiffParent   = AliasableCommand{"diff-parent"}   //nolint:gochecknoglobals
	AliasableCommandHack         = AliasableCommand{"hack"}          //nolint:gochecknoglobals
	AliasableCommandKill         = AliasableCommand{"kill"}          //nolint:gochecknoglobals
	AliasableCommandPrepend      = AliasableCommand{"prepend"}       //nolint:gochecknoglobals
	AliasableCommandPropose      = AliasableCommand{"propose"}       //nolint:gochecknoglobals
	AliasableCommandRenameBranch = AliasableCommand{"rename-branch"} //nolint:gochecknoglobals
	AliasableCommandRepo         = AliasableCommand{"repo"}          //nolint:gochecknoglobals
	AliasableCommandShip         = AliasableCommand{"ship"}          //nolint:gochecknoglobals
	AliasableCommandSync         = AliasableCommand{"sync"}          //nolint:gochecknoglobals
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

func LookupAliasableCommand(key string) *AliasableCommand {
	if !strings.HasPrefix(key, "alias.") {
		return nil
	}
	for _, aliasableCommand := range AliasableCommands() {
		aliasKey := aliasableCommand.Key()
		if key == aliasKey.name {
			return &aliasableCommand
		}
	}
	return nil
}
