package kit

import (
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/tools/go/packages"
	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"

	"sagikazarmark.dev/mga/internal/generate/kit/endpoint/endpointgen"
	"sagikazarmark.dev/mga/pkg/genutils"
)

type endpointOptions struct {
	headerFile string
	year       string

	paths  []string
	output string
}

// NewEndpointCommand returns a cobra command for generating an endpoint.
func NewEndpointCommand() *cobra.Command {
	var options endpointOptions

	cmd := &cobra.Command{
		Use:     "endpoint [flags] [paths]",
		Aliases: []string{"e"},
		Short:   "Generate Go kit endpoints from service interfaces",
		Long: `This command generates type safe Go kit endpoint structs.

Service interfaces look like the following:

	type Service interface {
		Call(ctx context.Context, req interface{}) (interface{}, error)

		// ... other calls
	}

where request and response types are any structures in the package.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceErrors = true
			cmd.SilenceUsage = true

			options.paths = args

			return runEndpoint(options)
		},
	}

	flags := cmd.Flags()

	flags.StringVar(&options.output, "output", "subpkg:suffix=driver", "output rule")
	flags.StringVar(&options.headerFile, "header-file", "", "header text (e.g. license) to prepend to generated files")
	flags.StringVar(&options.year, "year", "", "copyright year")

	return cmd
}

func runEndpoint(options endpointOptions) error {
	var generator genall.Generator = endpointgen.Generator{
		HeaderFile: options.headerFile,
		Year:       options.year,
	}

	generators := genall.Generators{&generator}

	if len(options.paths) == 0 {
		options.paths = []string{"."}
	}

	runtime, err := forRoots(generators, options.paths...)
	// runtime, err := generators.ForRoots(options.paths...)
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

// copied from genall package to override package loader configuration.
//
// required for supporting various types (basic type aliases, imports from other packages).
func forRoots(g genall.Generators, rootPaths ...string) (*genall.Runtime, error) {
	roots, err := loader.LoadRootsWithConfig(
		&packages.Config{
			Mode: packages.NeedDeps | packages.NeedTypes,
		},
		rootPaths...,
	)
	if err != nil {
		return nil, err
	}
	rt := &genall.Runtime{
		Generators: g,
		GenerationContext: genall.GenerationContext{
			Collector: &markers.Collector{
				Registry: &markers.Registry{},
			},
			Roots:     roots,
			InputRule: genall.InputFromFileSystem,
			Checker:   &loader.TypeChecker{},
		},
		OutputRules: genall.OutputRules{Default: genall.OutputToNothing},
	}
	if err := rt.Generators.RegisterMarkers(rt.Collector.Registry); err != nil {
		return nil, err
	}

	return rt, nil
}
