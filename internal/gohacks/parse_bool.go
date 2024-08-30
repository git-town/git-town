package gohacks

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/git-town/git-town/v16/internal/messages"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

func ParseBool(text, source string) (Option[bool], error) {
	switch strings.ToLower(text) {
	case "":
		return None[bool](), nil
	case "yes", "on", "enable", "enabled":
		return Some(true), nil
	case "no", "off", "disable", "disabled":
		return Some(false), nil
	}
	parsed, err := strconv.ParseBool(text)
	if err != nil {
		return None[bool](), fmt.Errorf(messages.ValueInvalid, source, text)
	}
	return Some(parsed), nil
}
