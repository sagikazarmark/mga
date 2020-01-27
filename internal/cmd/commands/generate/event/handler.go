package event

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"sagikazarmark.dev/mga/internal/generate/event/handler"
	"sagikazarmark.dev/mga/internal/generate/gentypes"
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

	event, err := handler.Parse(indir, options.event)
	if err != nil {
		return err
	}

	var packageRef gentypes.PackageRef
	var absOutDir string

	if options.outdir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		options.outdir = filepath.Base(cwd) + "gen"
		packageRef.Name = filepath.Base(options.outdir)

		absOut, err := filepath.Abs(options.outdir)
		if err != nil {
			return err
		}

		absOutDir = absOut
	} else {
		absOut, err := filepath.Abs(options.outdir)
		if err != nil {
			return err
		}
		absOutDir = absOut

		packageRef.Name = filepath.Base(absOut)

		absIn, err := filepath.Abs(indir)
		if err != nil {
			return err
		}

		if absIn == absOut { // When the input and the output directories are the same
			packageRef = event.Package
		}
	}

	err = os.MkdirAll(absOutDir, 0755)
	if err != nil {
		return err
	}

	file := handler.File{
		File: gentypes.File{
			Package:    packageRef,
			HeaderText: "",
		},
		EventHandlers: []handler.EventHandler{
			{
				Name:  event.Name,
				Event: event,
			},
		},
	}

	resFile := filepath.Join(absOutDir, fmt.Sprintf("%s_event_handler_gen.go", event.Name))

	fmt.Printf("Generating event handler for %s in %s\n", event.Name, resFile)

	res, err := handler.Generate(file)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(resFile, res, 0644)
	if err != nil {
		return err
	}

	return nil
}
