package program

import (
	"github.com/git-town/git-town/v9/src/gohacks"
	"github.com/git-town/git-town/v9/src/vm/opcode"
)

func IsCheckoutStep(step opcode.Opcode) bool {
	typeName := gohacks.TypeName(step)
	return typeName == "Checkout" || typeName == "CheckoutIfExists"
}
