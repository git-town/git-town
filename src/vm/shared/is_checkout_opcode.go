package shared

import (
	"github.com/git-town/git-town/v12/src/gohacks"
)

func IsCheckoutOpcode(opcode Opcode) bool {
	typeName := gohacks.TypeName(opcode)
	return typeName == "Checkout" || typeName == "CheckoutIfExists"
}
