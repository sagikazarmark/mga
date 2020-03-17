package scaffold

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"sagikazarmark.dev/mga/internal/scaffold/service"
)

type serviceOptions struct {
	path string

	force bool
}

// NewServiceCommand returns a cobra command for scaffolding a service interface.
func NewServiceCommand() *cobra.Command {
	var options serviceOptions

	cmd := &cobra.Command{
		Use:   "service [flags] PATH",
		Short: "Scaffold a new service interface",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceErrors = true
			cmd.SilenceUsage = true

			options.path = args[0]

			return runService(options)
		},
	}

	flags := cmd.Flags()

	flags.BoolVar(&options.force, "force", false, "overwrite existing file (BE CAREFUL!!!)")

	return cmd
}

func runService(options serviceOptions) error {
	path := filepath.Clean(options.path)
	pkg := filepath.Base(path)

	if pkg == "." || pkg == string(filepath.Separator) {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		pkg = filepath.Base(dir)
	}

	svcPath := filepath.Join(path, "service.go")

	if fileExists(svcPath) && !options.force {
		return errors.New("service already exists: use --force to overwrite it")
	}

	content, err := service.Scaffold(pkg)
	if err != nil {
		return err
	}

	file, err := os.Create(svcPath)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.Write(content)
	if err != nil {
		return err
	}

	return nil
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}
