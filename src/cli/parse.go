package cli

import (
	"strconv"
	"strings"
)

func ParseBool(text string) (bool, error) {
	text = strings.ToLower(text)
	switch text {
	case "yes", "on":
		return true, nil
	case "no", "off":
		return false, nil
	}
	return strconv.ParseBool(text)
}
