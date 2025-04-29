package gohacks

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/git-town/git-town/v19/internal/messages"
	. "github.com/git-town/git-town/v19/pkg/prelude"
)

func ParseBool(text, source string) (bool, error) {
	switch strings.ToLower(text) {
	case "":
		return false, fmt.Errorf(messages.ValueInvalid, source, text)
	case "yes", "on", "enable", "enabled":
		return true, nil
	case "no", "off", "disable", "disabled":
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
