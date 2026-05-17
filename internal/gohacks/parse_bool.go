package gohacks

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/git-town/git-town/v23/internal/gohacks/stringss"
	"github.com/git-town/git-town/v23/internal/messages"
	. "github.com/git-town/git-town/v23/pkg/prelude"
)

func ParseBool[T ~bool](text stringss.Trimmed, source string) (T, error) { //nolint:ireturn
	switch strings.ToLower(text.String()) {
	case "":
		return false, fmt.Errorf(messages.ValueInvalid, source, text)
	case "yes", "y", "on", "enable", "enabled":
		return true, nil
	case "no", "n", "off", "disable", "disabled":
		return false, nil
	}
	parsed, err := strconv.ParseBool(text.String())
	if err != nil {
		return false, fmt.Errorf(messages.ValueInvalid, source, text)
	}
	return T(parsed), nil
}

func ParseBoolOpt[T ~bool](text stringss.Trimmed, source string) (Option[T], error) {
	if text == "" {
		return None[T](), nil
	}
	parsed, err := ParseBool[T](text, source)
	return Some(parsed), err
}

func StrOpt2BoolOpt[T ~bool](textOpt Option[string], source string) (Option[T], error) {
	text, has := textOpt.Get()
	if !has {
		return None[T](), nil
	}
	parsed, err := ParseBool[T](text, source)
	return Some(parsed), err
}
