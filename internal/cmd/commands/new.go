package commands

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/fatih/color"
	"github.com/hashicorp/go-getter"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb/v4"
	"github.com/vbauerster/mpb/v4/decor"
)

// nolint: lll
const templateURL = "https://github.com/sagikazarmark/modern-go-application/archive/master.zip//modern-go-application-master"

type newOptions struct {
	destination string
	noInit      bool

	noProgress bool
}

// NewNewCommand returns a cobra command for creating a new project from the Modern Go Application template.
func NewNewCommand() *cobra.Command {
	var options newOptions

	cmd := &cobra.Command{
		Use:     "new [flags] DESTINATION",
		Aliases: []string{"n"},
		Short:   "Create a new project based on the Modern Go Application template.",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceErrors = true
			cmd.SilenceUsage = true

			options.destination = args[0]

			return runNew(cmd, options)
		},
	}

	flags := cmd.Flags()

	flags.BoolVar(&options.noInit, "no-init", false, "Do not initialize project after creation")
	flags.BoolVar(&options.noProgress, "no-progress", false, "Do not show download progress")

	return cmd
}

func runNew(cmd *cobra.Command, options newOptions) error {
	dest, err := filepath.Abs(options.destination)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	err = ensureEmptyDir(dest)
	if err != nil {
		return err
	}

	getterOptions := make([]getter.ClientOption, 0)

	cmd.Println(color.BlueString("Creating a new project under %q", dest))

	pb := newProgressBar(cmd)
	if !options.noProgress {
		cmd.Println()
		cmd.Println(color.BlueString("Downloading template files"))

		getterOptions = append(getterOptions, getter.WithProgress(pb))
	}

	err = getter.Get(dest, templateURL, getterOptions...)
	pb.progress.Wait()
	if err != nil {
		return err
	}

	if !options.noProgress {
		cmd.Println()
	}

	if !options.noInit && false {
		err := os.Chdir(dest)
		if err != nil {
			return fmt.Errorf("failed to enter project directory: %w", err)
		}

		cmd.Println(color.BlueString("Initializing project"))

		_, err = runSurvey()
		if err == terminal.InterruptErr {
			cmd.Println(color.YellowString(`Interrupted! Run "mga init" in your project to finish initializing it.`))

			return nil
		} else if err != nil {
			return err
		}
	} else { // nolint: staticcheck
		// nolint: lll
		// cmd.Println(color.YellowString(`Project is not initialized yet! Run "mga init" in your project to finish initializing it.`))
	}

	cmd.Println(color.GreenString("Project successfully created!"))

	return nil
}

func ensureEmptyDir(path string) error {
	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) { // Path does not exist, moving on
		return nil
	}

	if err != nil {
		return fmt.Errorf("failed to stat destination path: %w", err)
	}

	if !fileInfo.IsDir() {
		return errors.New("destination path already exists and is not a directory")
	}

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open destination directory: %w", err)
	}

	_, err = file.Readdirnames(1)
	if err != io.EOF {
		return errors.New("destination path already exists and is not an empty directory")
	}

	// Path is an empty directory
	return nil
}

type progressBar struct {
	// lock everything below
	lock sync.Mutex

	progress *mpb.Progress
}

func newProgressBar(cmd *cobra.Command) *progressBar {
	return &progressBar{
		progress: mpb.New(mpb.WithOutput(cmd.OutOrStderr())),
	}
}

// TrackProgress instantiates a new progress bar that will
// display the progress of stream until closed.
// total can be 0.
func (cpb *progressBar) TrackProgress(_ string, _, totalSize int64, stream io.ReadCloser) io.ReadCloser {
	cpb.lock.Lock()
	defer cpb.lock.Unlock()

	if cpb.progress == nil {
		cpb.progress = mpb.New()
	}
	bar := cpb.progress.AddBar(
		totalSize,
		mpb.PrependDecorators(
			decor.CountersKibiByte("% 6.1f / % 6.1f"),
		),
		mpb.AppendDecorators(
			decor.EwmaETA(decor.ET_STYLE_MMSS, float64(totalSize)/2048),
			decor.Name(" ] "),
			decor.AverageSpeed(decor.UnitKiB, "% .2f"),
		),
	)

	reader := bar.ProxyReader(stream)

	return &readCloser{
		Reader: reader,
		close: func() error {
			cpb.lock.Lock()
			defer cpb.lock.Unlock()

			return nil
		},
	}
}

type readCloser struct {
	io.Reader
	close func() error
}

func (c *readCloser) Close() error { return c.close() }
