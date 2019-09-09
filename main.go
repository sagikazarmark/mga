package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"sagikazarmark.dev/mga/internal/cmd/commands"
)

// Provisioned by ldflags
// nolint: gochecknoglobals
var (
	version    string
	commitHash string
	buildDate  string
)

const (
	// appName is an identifier-like name used anywhere this app needs to be identified.
	//
	// It identifies the service itself, the actual instance needs to be identified via environment
	// and other details.
	appName = "mga"
)

func main() {
	// rootCmd represents the base command when called without any subcommands
	rootCmd := &cobra.Command{
		Use:     appName,
		Short:   "CLI tool for Modern Go Application based apps",
		Version: version,
	}

	rootCmd.SetVersionTemplate(fmt.Sprintf("%s version %s (%s) built on %s\n", appName, version, commitHash, buildDate))

	commands.AddCommands(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)

		os.Exit(1)
	}
}
