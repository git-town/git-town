package configdomain

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v23/internal/messages"
	. "github.com/git-town/git-town/v23/pkg/prelude"
)

// Order describes the sort order for displaying branch lists.
type Order string

const (
	OrderAsc  Order = "asc"
	OrderDesc Order = "desc"
)

var OrderValues = []Order{
	OrderAsc,
	OrderDesc,
}

func (self Order) String() string {
	return string(self)
}

func ParseOrder(value string, source string) (Option[Order], error) {
	switch strings.ToLower(value) {
	case "a", "as", "asc":
		return Some(OrderAsc), nil
	case "d", "de", "des", "desc":
		return Some(OrderDesc), nil
	default:
		return None[Order](), fmt.Errorf(messages.OrderInvalid, source, value)
	}
}

func ParseOrderOpt(valueOpt Option[string], source string) (Option[Order], error) {
	if value, has := valueOpt.Get(); has {
		return ParseOrder(value, source)
	}
	return None[Order](), nil
}
