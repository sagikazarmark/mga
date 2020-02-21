package testify

import (
	"github.com/spf13/cobra"
)

// NewTestifyCommand returns a cobra command for `testify` subcommands.
func NewTestifyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "testify",
		Aliases: []string{"t"},
		Short:   "Generate testify code",
	}

	cmd.AddCommand(
		NewMockCommand(),
	)

	return cmd
}
