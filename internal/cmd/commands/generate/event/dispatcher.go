package event

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"sagikazarmark.dev/mga/internal/generate/event/dispatcher"
	"sagikazarmark.dev/mga/internal/generate/gentypes"
)

type dispatcherOptions struct {
	baseInterface string
	outdir        string
}

// NewDispatcherCommand returns a cobra command for generating an event dispatcher.
func NewDispatcherCommand() *cobra.Command {
	var options dispatcherOptions

	cmd := &cobra.Command{
		Use:     "dispatcher [options] INTERFACE",
		Aliases: []string{"d", "disp"},
		Short:   "Generate an event dispatcher from a base interface",
		Long: `This command generates a type safe event dispatcher implementation with an underlying generic event bus.
The event bus itself is an interface generated alongside the dispatcher:

	type EventBus interface {
		Publish(ctx context.Context, event interface{}) error
	}

You can either implement this interface yourself or use an implementation that's already compatible with it
(for example Watermill: https://github.com/ThreeDotsLabs/watermill).

Base interfaces look like the following:

	type Events interface {
		Event(ctx context.Context, event Event) error

		// ... other events
	}

where Event is a simple data structure containing the event payload.
The context parameter and the error return value are both optional,
but interface methods cannot accept or return more or different parameters.
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceErrors = true
			cmd.SilenceUsage = true

			options.baseInterface = args[0]

			return runDispatcher(options)
		},
	}

	flags := cmd.Flags()

	flags.StringVar(&options.outdir, "outdir", "", "output directory (default: $PWD/currdir+'gen', eg. module/modulegen)")

	return cmd
}

func runDispatcher(options dispatcherOptions) error {
	indir := "."

	events, err := dispatcher.Parse(indir, options.baseInterface)
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
			packageRef = events.Package
		}
	}

	err = os.MkdirAll(absOutDir, 0755)
	if err != nil {
		return err
	}

	resFile := filepath.Join(absOutDir, fmt.Sprintf("%s_event_dispatcher_gen.go", events.Name))

	file := dispatcher.File{
		File: gentypes.File{
			Package:    packageRef,
			HeaderText: "",
		},
		EventDispatchers: []dispatcher.EventDispatcher{
			{
				Name:              cleanEventDispatcherName(events.Name),
				DispatcherMethods: events.Methods,
			},
		},
	}

	fmt.Printf("Generating event dispatcher for %s in %s\n", events.Name, resFile)

	res, err := dispatcher.Generate(file)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(resFile, res, 0644)
	if err != nil {
		return err
	}

	return nil
}

func cleanEventDispatcherName(name string) string {
	name = strings.TrimSuffix(name, "Events")
	name = strings.TrimSuffix(name, "EventBus")
	name = strings.TrimSuffix(name, "EventDispatcher")

	return name
}
