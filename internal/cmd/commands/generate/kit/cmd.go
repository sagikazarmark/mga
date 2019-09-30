package kit

import (
	"github.com/spf13/cobra"
)

// NewKitCommand returns a cobra command for `kit` subcommands.
func NewKitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "kit",
		Aliases: []string{"k"},
		Short:   "Generate go-kit code",
	}

	cmd.AddCommand(
		NewEndpointCommand(),
	)

	return cmd
}
