package program

import (
	"github.com/git-town/git-town/v9/src/gohacks"
	"github.com/git-town/git-town/v9/src/step"
)

func IsCheckoutStep(step step.Step) bool {
	typeName := gohacks.TypeName(step)
	return typeName == "Checkout" || typeName == "CheckoutIfExists"
}
