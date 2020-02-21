package generate

import (
	"github.com/spf13/cobra"

	"sagikazarmark.dev/mga/internal/cmd/commands/generate/event"
	"sagikazarmark.dev/mga/internal/cmd/commands/generate/kit"
	"sagikazarmark.dev/mga/internal/cmd/commands/generate/testify"
)

// NewGenerateCommand returns a cobra command for `generate` subcommands.
func NewGenerateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate",
		Aliases: []string{"g", "gen"},
		Short:   "Generate code",
	}

	cmd.AddCommand(
		event.NewEventsCommand(),
		kit.NewKitCommand(),
		testify.NewTestifyCommand(),
		NewMockeryCommand(),
	)

	return cmd
}
