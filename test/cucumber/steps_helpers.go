package cucumber

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/config/configdomain"
)

func checkString(haveNullable *string, want string, name string) error {
	var have string
	if haveNullable == nil {
		have = ""
	} else {
		have = *haveNullable
	}
	if have != want {
		return fmt.Errorf("unexpected value for key %q: want %q have %q", configdomain.KeyAliasAppend, want, have)
	}
	return nil
}
