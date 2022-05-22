package commands

import (
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/gobuffalo/here"
)

type initSurvey struct {
	ModuleName      string
	ProjectName     string
	BinaryName      string
	AppName         string
	FriendlyAppName string
}

func runSurvey() (initSurvey, error) {
	var results initSurvey

	info, err := here.Current()
	if err != nil {
		return results, fmt.Errorf("could not get module information: %w", err)
	}

	if !info.Module.Main {
		return results, errors.New("this project does not seem to be a Go application")
	}

	err = survey.AskOne(
		&survey.Input{Message: "Go module name?", Default: info.Module.Path},
		&results.ModuleName,
		survey.WithValidator(survey.Required),
	)
	if err != nil {
		return results, err
	}

	err = survey.AskOne(
		&survey.Input{
			Message: "Project name?",
			Default: path.Base(results.ModuleName),
			Help:    "A directory-like name for your project (used in configurations, etc)",
		},
		&results.ProjectName,
		survey.WithValidator(survey.Required),
	)
	if err != nil {
		return results, err
	}

	err = survey.AskOne(
		&survey.Input{
			Message: "Binary name?",
			Default: results.ProjectName,
			Help:    "Name of the main binary",
		},
		&results.BinaryName,
		survey.WithValidator(survey.Required),
	)
	if err != nil {
		return results, err
	}

	err = survey.AskOne(
		&survey.Input{
			Message: "Application name?",
			Default: results.BinaryName,
			Help:    "Application name (if it differs from the binary name)",
		},
		&results.AppName,
		survey.WithValidator(survey.Required),
	)
	if err != nil {
		return results, err
	}

	var friendlyAppName string
	err = survey.AskOne(
		&survey.Input{
			Message: "Friendly application name?",
			Default: strings.Title(strings.ReplaceAll(results.AppName, "-", " ")), // nolint
			Help:    "User friendly application name (appearing eg. in logs)",
		},
		&friendlyAppName,
		survey.WithValidator(survey.Required),
	)
	if err != nil {
		return results, err
	}

	return results, nil
}
