package flags

import (
	"cmp"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/spf13/cobra"
)

const orderLong = "order"

// type-safe access to the CLI arguments of type configdomain.Order
func Order() (AddFunc, ReadOrderFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().StringP(orderLong, "o", "", "sort order for branch list (asc or desc)")
	}
	readFlag := func(cmd *cobra.Command) (Option[configdomain.Order], error) {
		text, errFlag := cmd.Flags().GetString(orderLong)
		order, errParse := configdomain.ParseOrder(text, "--order flag")
		return order, cmp.Or(errFlag, errParse)
	}
	return addFlag, readFlag
}

// ReadOrderFlagFunc is the type signature for the function that reads the "order" flag from the args to the given Cobra command.
type ReadOrderFlagFunc func(*cobra.Command) (Option[configdomain.Order], error)
