package prompt

import (
	"runtime"

	"github.com/Originate/exit"
	survey "gopkg.in/AlecAivazis/survey.v1"
	surveyCore "gopkg.in/AlecAivazis/survey.v1/core"
)

type askForBranchOptions struct {
	branchNames       []string
	prompt            string
	defaultBranchName string
}

type askForBranchesOptions struct {
	branchNames        []string
	prompt             string
	defaultBranchNames []string
}

func askForBranch(opts askForBranchOptions) string {
	result := ""
	if runtime.GOOS == "windows" {
		surveyCore.SelectFocusIcon = ">"
	}
	prompt := &survey.Select{
		Message: opts.prompt,
		Options: opts.branchNames,
		Default: opts.defaultBranchName,
	}
	err := survey.AskOne(prompt, &result, nil)
	exit.If(err)
	return result
}

func askForBranches(opts askForBranchesOptions) []string {
	result := []string{}
	if runtime.GOOS == "windows" {
		surveyCore.MarkedOptionIcon = "[x]"
		surveyCore.UnmarkedOptionIcon = "[ ]"
	}
	prompt := &survey.MultiSelect{
		Message: opts.prompt,
		Options: opts.branchNames,
		Default: opts.defaultBranchNames,
	}
	err := survey.AskOne(prompt, &result, nil)
	exit.If(err)
	return result
}
