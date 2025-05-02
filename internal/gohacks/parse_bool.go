package gohacks

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/git-town/git-town/v20/internal/messages"
	. "github.com/git-town/git-town/v20/pkg/prelude"
)

func ParseBool(text, source string) (bool, error) {
	switch strings.ToLower(text) {
	case "":
		return false, fmt.Errorf(messages.ValueInvalid, source, text)
	case "yes", "y", "on", "enable", "enabled":
		return true, nil
	case "no", "n", "off", "disable", "disabled":
		return false, nil
	}
	result, err := strconv.ParseBool(text)
	if err != nil {
		return false, fmt.Errorf(messages.ValueInvalid, source, text)
	}
	return result, nil
}

func ParseBoolOpt(text, source string) (Option[bool], error) {
	if text == "" {
		return None[bool](), nil
	}
	result, err := ParseBool(text, source)
	return Some(result), err
}
