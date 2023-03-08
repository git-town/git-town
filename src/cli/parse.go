package cli

import (
	"strconv"
	"strings"
)

func ParseBool(text string) (bool, error) {
	text = strings.ToLower(text)
	if text == "yes" || text == "on" {
		return true, nil
	}
	if text == "no" || text == "off" {
		return false, nil
	}
	return strconv.ParseBool(text)
}
