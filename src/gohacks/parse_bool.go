package gohacks

import (
	"strconv"
	"strings"
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
