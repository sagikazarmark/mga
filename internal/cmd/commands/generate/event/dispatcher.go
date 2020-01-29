package event

import (
	"os"

	"github.com/spf13/cobra"
	"sigs.k8s.io/controller-tools/pkg/genall"

	"sagikazarmark.dev/mga/internal/generate/event/dispatcher/dispatchergen"
	"sagikazarmark.dev/mga/pkg/genutils"
)

type dispatcherOptions struct {
	headerFile string
	year       string

	paths []string
}

// NewDispatcherCommand returns a cobra command for generating an event dispatcher.
func NewDispatcherCommand() *cobra.Command {
	var options dispatcherOptions

	cmd := &cobra.Command{
		Use:     "dispatcher [options] [paths]",
		Aliases: []string{"d", "disp"},
		Short:   "Generate a event dispatcher implementations from base interfaces",
		Long: `This command generates type safe event dispatcher implementations with an underlying generic event bus.
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

			options.paths = args

			return runDispatcher(options)
		},
	}

	flags := cmd.Flags()

	flags.StringVar(&options.headerFile, "header-file", "", "header text (e.g. license) to prepend to generated files")
	flags.StringVar(&options.year, "year", "", "copyright year")

	return cmd
}

func runDispatcher(options dispatcherOptions) error {
	var generator genall.Generator = dispatchergen.Generator{
		HeaderFile: options.headerFile,
		Year:       options.year,
	}

	generators := genall.Generators{&generator}

	if len(options.paths) == 0 {
		options.paths = []string{"."}
	}

	runtime, err := generators.ForRoots(options.paths...)
	if err != nil {
		return err
	}

	runtime.OutputRules.Default = genutils.OutputArtifacts{}

	if hadErrs := runtime.Run(); hadErrs {
		os.Exit(1)
	}

	return nil
}
