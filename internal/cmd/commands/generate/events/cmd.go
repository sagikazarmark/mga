package events

import (
	"github.com/spf13/cobra"
)

// NewEventsCommand returns a cobra command for `events` subcommands.
func NewEventsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "events",
		Aliases: []string{"e", "ev", "event"},
		Short:   "Generate event related code",
	}

	cmd.AddCommand(
		NewDispatcherCommand(),
	)

	return cmd
}
