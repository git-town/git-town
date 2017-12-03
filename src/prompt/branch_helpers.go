package prompt

import (
	"github.com/Originate/exit"
	survey "gopkg.in/AlecAivazis/survey.v1"
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
	prompt := &survey.MultiSelect{
		Message: opts.prompt,
		Options: opts.branchNames,
		Default: opts.defaultBranchNames,
	}
	err := survey.AskOne(prompt, &result, nil)
	exit.If(err)
	return result
}
