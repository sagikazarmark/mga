package event

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"sagikazarmark.dev/mga/internal/generate/event/handler"
)

type handlerOptions struct {
	event  string
	outdir string
}

// NewHandlerCommand returns a cobra command for generating an event handler.
func NewHandlerCommand() *cobra.Command {
	var options handlerOptions

	cmd := &cobra.Command{
		Use:     "handler [options] INTERFACE",
		Aliases: []string{"h"},
		Short:   "Generate an event handler for an event",
		Long: `This command generates a type safe event handler implementation for an event.
An event can be any plain, exported struct:

	type Event struct {
		ID string
	}

The generated handler is compatible with Watermill (https://github.com/ThreeDotsLabs/watermill).
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceErrors = true
			cmd.SilenceUsage = true

			options.event = args[0]

			return runHandler(options)
		},
	}

	flags := cmd.Flags()

	flags.StringVar(&options.outdir, "outdir", "", "output directory (default: $PWD/currdir+'gen', eg. module/modulegen)")

	return cmd
}

func runHandler(options handlerOptions) error {
	indir := "."

	spec, err := handler.Parse(indir, options.event)
	if err != nil {
		return err
	}

	var outpkg string

	if options.outdir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		options.outdir = filepath.Base(cwd) + "gen"
		outpkg = filepath.Base(options.outdir)
	} else {
		absOut, err := filepath.Abs(options.outdir)
		if err != nil {
			return err
		}

		outpkg = filepath.Base(absOut)

		absIn, err := filepath.Abs(indir)
		if err != nil {
			return err
		}

		if absIn == absOut { // When the input and the output directories are the same
			outpkg = spec.Package.Path
		}
	}

	err = os.MkdirAll(options.outdir, 0755)
	if err != nil {
		return err
	}

	res, err := handler.Generate(outpkg, spec)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(
		filepath.Join(
			options.outdir,
			fmt.Sprintf("%s_event_handler_gen.go", spec.Name),
		),
		[]byte(res),
		0644,
	)
	if err != nil {
		return err
	}

	return nil
}
