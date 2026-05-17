package configdomain

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v23/internal/gohacks/stringss"
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

func ParseOrder(value stringss.Trimmed, source string) (Option[Order], error) {
	switch strings.ToLower(value.String()) {
	case "":
		return None[Order](), nil
	case "a", "as", "asc":
		return Some(OrderAsc), nil
	case "d", "de", "des", "desc":
		return Some(OrderDesc), nil
	default:
		return None[Order](), fmt.Errorf(messages.OrderInvalid, source, value)
	}
}
