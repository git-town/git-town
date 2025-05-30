package configdomain

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

const (
	ShipStrategyAPI         ShipStrategy = "api"          // shipping via the forge API
	ShipStrategyAlwaysMerge ShipStrategy = "always-merge" // shipping by doing a local merge commit (merge --no-ff)
	ShipStrategyFastForward ShipStrategy = "fast-forward" // shipping by doing a local fast-forward
	ShipStrategySquashMerge ShipStrategy = "squash-merge" // shipping by doing a local squash-merge
)

type ShipStrategy string

func (self ShipStrategy) String() string {
	return string(self)
}

func ParseShipStrategy(text string) (Option[ShipStrategy], error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return None[ShipStrategy](), nil
	}
	text = strings.ToLower(text)
	for _, shipStrategy := range ShipStrategies() {
		if shipStrategy.String() == text {
			return Some(shipStrategy), nil
		}
	}
	return None[ShipStrategy](), fmt.Errorf(messages.ConfigShipStrategyUnknown, text)
}

func ShipStrategies() []ShipStrategy {
	return []ShipStrategy{
		ShipStrategyAPI,
		ShipStrategyAlwaysMerge,
		ShipStrategyFastForward,
		ShipStrategySquashMerge,
	}
}
