package generate

import (
	"fmt"
	"os"
	"regexp"
	"runtime/pprof"
	"strings"

	"emperror.dev/errors"
	"github.com/spf13/cobra"
	"github.com/vektra/mockery/mockery"
)

type mockeryOptions struct {
	fName      string
	fPrint     bool
	fOutput    string
	fOutpkg    string
	fDir       string
	fRecursive bool
	fAll       bool
	fIP        bool
	fTO        bool
	fCase      string
	fNote      string
	fProfile   string
	fkeepTree  bool
	buildTags  string
}

// NewMockeryCommand returns a cobra command for generating a mock using mockery.
func NewMockeryCommand() *cobra.Command {
	var options mockeryOptions

	cmd := &cobra.Command{
		Use:   "mockery",
		Short: "Generate a mock from an interface using Mockery",
		Long: `This command is a drop-in replacement for Mockery.

It uses the original code base under https://github.com/vektra/mockery

The command accepts the same arguments as the original executable.
`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceErrors = true
			cmd.SilenceUsage = true

			return runMockery(options)
		},
	}

	flags := cmd.Flags()

	flags.StringVar(&options.fName, "name", "", "name or matching regular expression of interface to generate mock for")
	flags.BoolVar(&options.fPrint, "print", false, "print the generated mock to stdout")
	flags.StringVar(&options.fOutput, "output", "./mocks", "directory to write mocks to")
	flags.StringVar(&options.fOutpkg, "outpkg", "mocks", "name of generated package")
	flags.StringVar(&options.fDir, "dir", ".", "directory to search for interfaces")
	flags.BoolVar(&options.fRecursive, "recursive", false, "recurse search into sub-directories")
	flags.BoolVar(&options.fAll, "all", false, "generates mocks for all found interfaces in all sub-directories")
	flags.BoolVar(&options.fIP, "inpkg", false, "generate a mock that goes inside the original package")
	flags.BoolVar(&options.fTO, "testonly", false, "generate a mock in a _test.go file")
	flags.StringVar(&options.fCase, "case", "camel", "name the mocked file using casing convention [camel, snake, underscore]") // nolint: lll
	flags.StringVar(&options.fNote, "note", "", "comment to insert into prologue of each generated file")
	flags.StringVar(&options.fProfile, "cpuprofile", "", "write cpu profile to file")
	flags.BoolVar(&options.fkeepTree, "keeptree", false, "keep the tree structure of the original interface files into a different repository. Must be used with XX") // nolint: lll
	flags.StringVar(&options.buildTags, "tags", "", "space-separated list of additional build tags to use")

	return cmd
}

const regexMetadataChars = "\\.+*?()|[]{}^$"

func runMockery(options mockeryOptions) error {
	var recursive bool
	var filter *regexp.Regexp
	var err error
	var limitOne bool

	// nolint: gocritic
	if options.fName != "" && options.fAll {
		return errors.New("specify -name or -all, but not both")
	} else if options.fName != "" {
		recursive = options.fRecursive
		if strings.ContainsAny(options.fName, regexMetadataChars) {
			if filter, err = regexp.Compile(options.fName); err != nil {
				return errors.New("invalid regular expression provided to -name")
			}
		} else {
			filter = regexp.MustCompile(fmt.Sprintf("^%s$", options.fName))
			limitOne = true
		}
	} else if options.fAll {
		recursive = true
		filter = regexp.MustCompile(".*")
	} else {
		return errors.New("use -name to specify the name of the interface or -all for all interfaces found")
	}

	if options.fkeepTree {
		options.fIP = false
	}

	if options.fProfile != "" {
		f, err := os.Create(options.fProfile)
		if err != nil {
			return err
		}

		_ = pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	var osp mockery.OutputStreamProvider
	if options.fPrint {
		osp = &mockery.StdoutStreamProvider{}
	} else {
		osp = &mockery.FileOutputStreamProvider{
			BaseDir:                   options.fOutput,
			InPackage:                 options.fIP,
			TestOnly:                  options.fTO,
			Case:                      options.fCase,
			KeepTree:                  options.fkeepTree,
			KeepTreeOriginalDirectory: options.fDir,
		}
	}

	visitor := &mockery.GeneratorVisitor{
		InPackage:   options.fIP,
		Note:        options.fNote,
		Osp:         osp,
		PackageName: options.fOutpkg,
	}

	walker := mockery.Walker{
		BaseDir:   options.fDir,
		Recursive: recursive,
		Filter:    filter,
		LimitOne:  limitOne,
		BuildTags: strings.Split(options.buildTags, " "),
	}

	generated := walker.Walk(visitor)

	if options.fName != "" && !generated {
		return errors.Errorf("unable to find %s in any go files under this path\n", options.fName)
	}

	return nil
}
