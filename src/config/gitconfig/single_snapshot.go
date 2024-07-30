package gitconfig

import "github.com/git-town/git-town/v14/src/config/configdomain"

// SingleSnapshot contains all of the local or global Git metadata config settings.
type SingleSnapshot map[configdomain.Key]string
