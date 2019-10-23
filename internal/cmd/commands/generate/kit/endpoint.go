package kit

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"sagikazarmark.dev/mga/internal/generate/kit/endpoint"
)

type endpointOptions struct {
	serviceInterface string
	outdir           string
	outfile          string
	withOc           bool
}

// NewEndpointCommand returns a cobra command for generating an endpoint.
func NewEndpointCommand() *cobra.Command {
	var options endpointOptions

	cmd := &cobra.Command{
		Use:     "endpoint [options] INTERFACE",
		Aliases: []string{"e"},
		Short:   "Generate an endpoint from a service interface",
		Long: `This command generates a type safe endpoint struct.

Service interfaces look like the following:

	type Service interface {
		Call(ctx context.Context, req interface{}) (interface{}, error)

		// ... other calls
	}

where request and response types are any structures in the package.
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceErrors = true
			cmd.SilenceUsage = true

			options.serviceInterface = args[0]

			return runEndpoint(options)
		},
	}

	flags := cmd.Flags()

	flags.StringVar(&options.outdir, "outdir", "", "output directory (default: $PWD/currdir+'gen', eg. module/modulegen)")
	flags.StringVar(&options.outfile, "outfile", "endpoint_gen.go", "output file within the output directory")
	flags.BoolVar(&options.withOc, "with-oc", false, "generate OpenCensus tracing middleware")

	return cmd
}

func runEndpoint(options endpointOptions) error {
	indir := "."

	spec, err := endpoint.Parse(indir, options.serviceInterface)
	if err != nil {
		return err
	}

	var outpkg string
	var absOutDir string

	if options.outdir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		options.outdir = filepath.Base(cwd) + "gen"
		outpkg = filepath.Base(options.outdir)

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

		outpkg = filepath.Base(absOut)

		absIn, err := filepath.Abs(indir)
		if err != nil {
			return err
		}

		if absIn == absOut { // When the input and the output directories are the same
			outpkg = spec.Package.Path
		}
	}

	if options.outfile == "" {
		options.outfile = "endpoint_gen.go"
	}

	options.outfile = filepath.Base(options.outfile)

	err = os.MkdirAll(absOutDir, 0755)
	if err != nil {
		return err
	}

	resFile := filepath.Join(absOutDir, options.outfile)

	fmt.Printf("Generating Go kit endpoints for %s in %s\n", spec.Name, resFile)

	res, err := endpoint.Generate(outpkg, spec, options.withOc)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(resFile, []byte(res), 0644)
	if err != nil {
		return err
	}

	return nil
}
