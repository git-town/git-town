package cli

import "strconv"

func ParseBool(text string) (bool, error) {
	if text == "yes" || text == "on" {
		return true, nil
	}
	if text == "no" || text == "off" {
		return false, nil
	}
	return strconv.ParseBool(text)
}
