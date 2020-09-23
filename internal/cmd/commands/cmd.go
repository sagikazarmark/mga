package commands

import (
	"github.com/spf13/cobra"

	"sagikazarmark.dev/mga/internal/cmd/commands/generate"
	"sagikazarmark.dev/mga/internal/cmd/commands/scaffold"
)

// AddCommands adds all the commands from cli/command to the root command.
func AddCommands(cmd *cobra.Command) {
	cmd.AddCommand(
		NewNewCommand(),
		generate.NewGenerateCommand(),
		scaffold.NewScaffoldCommand(),
	)
}
