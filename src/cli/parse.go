package cli

import "strconv"

func ParseBool(text string) (bool, error) {
	switch text {
	case "yes":
		return true, nil
	case "on":
		return true, nil
	case "no":
		return false, nil
	case "off":
		return false, nil
	}
	return strconv.ParseBool(text)
}
