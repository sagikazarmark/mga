package scaffold

import (
	"github.com/spf13/cobra"
)

// NewScaffoldCommand returns a cobra command for `scaffold` subcommands.
func NewScaffoldCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "scaffold",
		Aliases: []string{"scaff", "create", "c"},
		Short:   "Scaffold code",
	}

	cmd.AddCommand(
		NewServiceCommand(),
	)

	return cmd
}
