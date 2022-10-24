package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"sagikazarmark.dev/mga/internal/cmd/commands"
)

// Provisioned by ldflags
// nolint: gochecknoglobals
var (
	version    string
	commitHash string
	commitDate string
)

const (
	// appName is an identifier-like name used anywhere this app needs to be identified.
	//
	// It identifies the service itself, the actual instance needs to be identified via environment
	// and other details.
	appName = "mga"
)

func main() {
	var noColor bool

	// rootCmd represents the base command when called without any subcommands
	rootCmd := &cobra.Command{
		Use:     appName,
		Short:   "CLI tool for Modern Go Application based apps",
		Version: version,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			color.NoColor = noColor
		},
	}

	rootCmd.SetVersionTemplate(fmt.Sprintf("%s version %s (%s) on %s\n", appName, version, commitHash, commitDate))

	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable colorized output")

	commands.AddCommands(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		rootCmd.PrintErrln(color.RedString(err.Error()))

		os.Exit(1)
	}
}
