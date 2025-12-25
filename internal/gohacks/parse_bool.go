package gohacks

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

func ParseBool[T ~bool](text, source string) (T, error) { //nolint:ireturn
	switch strings.ToLower(text) {
	case "":
		return false, fmt.Errorf(messages.ValueInvalid, source, text)
	case "yes", "y", "on", "enable", "enabled":
		return true, nil
	case "no", "n", "off", "disable", "disabled":
		return false, nil
	}
	parsed, err := strconv.ParseBool(text)
	if err != nil {
		return false, fmt.Errorf(messages.ValueInvalid, source, text)
	}
	return T(parsed), nil
}

func ParseBoolOpt[T ~bool](text, source string) (Option[T], error) {
	if text == "" {
		return None[T](), nil
	}
	parsed, err := ParseBool[T](text, source)
	return Some(parsed), err
}
