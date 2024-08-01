package gohacks

import (
	"strconv"
	"strings"

	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
)

func ParseBool(text string) (bool, error) {
	switch strings.ToLower(text) {
	case "yes", "on", "enable", "enabled":
		return true, nil
	case "no", "off", "disable", "disabled":
		return false, nil
	}
	return strconv.ParseBool(text)
}

func ParseBoolOpt(text string) (Option[bool], error) {
	if text == "" {
		return None[bool](), nil
	}
	parsed, err := ParseBool(text)
	return Some(parsed), err
}
