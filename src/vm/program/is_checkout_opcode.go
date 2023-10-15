package program

import (
	"github.com/git-town/git-town/v9/src/gohacks"
	"github.com/git-town/git-town/v9/src/vm/shared"
)

func IsCheckoutOpcode(opcode shared.Opcode) bool {
	typeName := gohacks.TypeName(opcode)
	return typeName == "Checkout" || typeName == "CheckoutIfExists"
}
