package cli

import (
	"strconv"
	"strings"
)

func ParseBool(text string) (bool, error) {
	switch strings.ToLower(text) {
	case "yes", "on":
		return true, nil
	case "no", "off":
		return false, nil
	}
	return strconv.ParseBool(text)
}
