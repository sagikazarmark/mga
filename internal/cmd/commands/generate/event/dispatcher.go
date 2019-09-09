package event

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"sagikazarmark.dev/mga/internal/generate/event/dispatcher"
)

type dispatcherOptions struct {
	from   string
	outdir string
}

// NewDispatcherCommand returns a cobra command for generating an event dispatcher.
func NewDispatcherCommand() *cobra.Command {
	var options dispatcherOptions

	cmd := &cobra.Command{
		Use:     "dispatcher",
		Aliases: []string{"d", "disp"},
		Short:   "Generate an event dispatcher from a base interface",
		Long: `This command generates a type safe event dispatcher implementation with an underlying, generic event bus.
The event bus itself is an interface:

	type EventBus interface {
		Publish(ctx context.Context, event interface{}) error
	}

You can either implement this interface yourself or use an implementation that's already compatible with it.
For example, Watermill (https://github.com/ThreeDotsLabs/watermill) already has an event bus that's compatible.

Base interfaces look like the following:

	type Events interface {
		Event(ctx context.Context, event Event) error

		// ... other events
	}

where Event is a simple data structure containing the event payload.
The context parameter and the error return value are both optional,
but interface methods cannot accept or return more or different parameters.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceErrors = true
			cmd.SilenceUsage = true

			return runDispatcher(options)
		},
	}

	flags := cmd.Flags()

	flags.StringVar(&options.from, "from", "Events", "base event dispatcher interface")
	flags.StringVar(&options.outdir, "outdir", "", "output directory (default: $PWD/currdir+'gen', eg. module/modulegen)")

	return cmd
}

func runDispatcher(options dispatcherOptions) error {
	indir := "."

	spec, err := dispatcher.Parse(indir, options.from)
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

	res, err := dispatcher.Generate(outpkg, spec)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(options.outdir, fmt.Sprintf("%s_event_dispatcher.go", spec.Name)), []byte(res), 0644)
	if err != nil {
		return err
	}

	return nil
}
