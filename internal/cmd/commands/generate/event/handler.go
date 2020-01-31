package event

import (
	"os"

	"github.com/spf13/cobra"
	"sigs.k8s.io/controller-tools/pkg/genall"

	"sagikazarmark.dev/mga/internal/generate/event/handler/handlergen"
	"sagikazarmark.dev/mga/pkg/genutils"
)

type handlerOptions struct {
	headerFile string
	year       string

	paths  []string
	output string
}

// NewHandlerCommand returns a cobra command for generating an event handler.
func NewHandlerCommand() *cobra.Command {
	var options handlerOptions

	cmd := &cobra.Command{
		Use:     "handler [flags] [paths]",
		Aliases: []string{"h"},
		Short:   "Generate event handlers for events",
		Long: `This command generates type safe event handler implementations for events.
An event can be any plain, exported struct:

	type Event struct {
		ID string
	}

The generated handler is compatible with Watermill (https://github.com/ThreeDotsLabs/watermill).
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceErrors = true
			cmd.SilenceUsage = true

			options.paths = args

			return runHandler(options)
		},
	}

	flags := cmd.Flags()

	flags.StringVar(&options.output, "output", "pkg", "output rule")
	flags.StringVar(&options.headerFile, "header-file", "", "header text (e.g. license) to prepend to generated files")
	flags.StringVar(&options.year, "year", "", "copyright year")

	return cmd
}

func runHandler(options handlerOptions) error {
	var generator genall.Generator = handlergen.Generator{
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

	outputRule, err := genutils.LookupOutput(options.output)
	if err != nil {
		return err
	}

	runtime.OutputRules.Default = outputRule

	if hadErrs := runtime.Run(); hadErrs {
		os.Exit(1)
	}

	return nil
}
