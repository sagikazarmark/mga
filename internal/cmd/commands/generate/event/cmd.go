package event

import (
	"github.com/spf13/cobra"
)

// NewEventsCommand returns a cobra command for `event` subcommands.
func NewEventsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "event",
		Aliases: []string{"e", "ev", "events"},
		Short:   "Generate event related code",
	}

	cmd.AddCommand(
		NewDispatcherCommand(),
		NewHandlerCommand(),
	)

	return cmd
}
