package kit

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"sagikazarmark.dev/mga/internal/generate/gentypes"
	"sagikazarmark.dev/mga/internal/generate/kit/endpoint"
)

type endpointOptions struct {
	serviceInterface string
	outdir           string
	outfile          string
	baseName         string
	withOc           bool
	ocRoot           string
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
	flags.StringVar(&options.baseName, "base-name", "", "add a base name to generated structs (default: none)")
	flags.BoolVar(&options.withOc, "with-oc", false, "generate OpenCensus tracing middleware")
	flags.StringVar(&options.ocRoot, "oc-root", "", "override the package name in the generated OC trace middleware")

	return cmd
}

func runEndpoint(options endpointOptions) error {
	indir := "."

	svc, err := endpoint.Parse(indir, options.serviceInterface)
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
			packageRef = svc.Package
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

	file := endpoint.File{
		File: gentypes.File{
			Package:    packageRef,
			HeaderText: "",
		},
		EndpointSets: []endpoint.EndpointSet{
			{
				Name:           strings.TrimSuffix(svc.Name, "Service"),
				Service:        svc.TypeRef,
				Endpoints:      nil,
				WithOpenCensus: options.withOc,
			},
		},
	}

	operationNameRoot := svc.Package.Name
	if options.ocRoot != "" {
		operationNameRoot = options.ocRoot
	}

	for _, method := range svc.Methods {
		var operationName string

		// if endpoint set name is empty, do not add it to the operation name
		if file.EndpointSets[0].Name == "" {
			operationName = fmt.Sprintf("%s.%s", operationNameRoot, method.Name)
		} else {
			operationName = fmt.Sprintf("%s.%s.%s", operationNameRoot, file.EndpointSets[0].Name, method.Name)
		}

		file.EndpointSets[0].Endpoints = append(
			file.EndpointSets[0].Endpoints,
			endpoint.Endpoint{
				Name:          method.Name,
				OperationName: operationName,
			},
		)
	}

	fmt.Printf("Generating Go kit endpoints for %s in %s\n", file.EndpointSets[0].Service.Name, resFile)

	res, err := endpoint.Generate(file)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(resFile, res, 0644)
	if err != nil {
		return err
	}

	return nil
}
