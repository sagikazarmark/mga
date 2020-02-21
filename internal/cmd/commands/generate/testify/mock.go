package testify

import (
	"os"

	"github.com/spf13/cobra"
	"sigs.k8s.io/controller-tools/pkg/genall"

	"sagikazarmark.dev/mga/internal/generate/testify/mock/mockgen"
	"sagikazarmark.dev/mga/pkg/genutils"
)

type mockOptions struct {
	headerFile string
	year       string
	testOnly   bool

	paths  []string
	output string
}

// NewMockCommand returns a cobra command for generating mocks.
func NewMockCommand() *cobra.Command {
	var options mockOptions

	cmd := &cobra.Command{
		Use:     "mock [flags] [paths]",
		Aliases: []string{"e"},
		Short:   "Generate Testify mocks from interfaces",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceErrors = true
			cmd.SilenceUsage = true

			options.paths = args

			return runMock(options)
		},
	}

	flags := cmd.Flags()

	flags.StringVar(&options.output, "output", "pkg", "output rule")
	flags.StringVar(&options.headerFile, "header-file", "", "header text (e.g. license) to prepend to generated files")
	flags.StringVar(&options.year, "year", "", "copyright year")
	flags.BoolVar(&options.testOnly, "test-only", false, "generate code in _test packages")

	return cmd
}

func runMock(options mockOptions) error {
	var generator genall.Generator = mockgen.Generator{
		HeaderFile: options.headerFile,
		Year:       options.year,
		TestOnly:   options.testOnly,
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
