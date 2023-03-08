package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/spf13/cobra"
)

func completionsCmd(rootCmd *cobra.Command) *cobra.Command {
	completionsNoDescFlag := false
	completionsCmd := cobra.Command{
		Use:   "completions [bash|zsh|fish|powershell]",
		Short: "Generates auto-completion for bash, zsh, fish, or PowerShell",
		Long: `Generates auto-completion for bash, zsh, fish, or PowerShell.
When set up, "git-town <TAB>" will auto-complete Git Town subcommands.

To load autocompletion for Bash, run this command:

	git-town completions bash | source

To load completions for each session, add the above line to your ~/.bashrc file.


To load autocompletion for Zsh, run this command:

	git-town completions zsh | source

To load completions for each session, add the above line to your ~/.zshrc file.


To load autocompletion for Fish, run this command:

	git-town completions fish | source

To load completions for each session, add the above line to your ~/.config/fish/config.fish file.


To load autocompletions for Powershell, run this command:

	git-town completions powershell | Out-String | Invoke-Expression

To load completions for each session, add the above line to your PowerShell profile.
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		Run: func(cmd *cobra.Command, args []string) {
			completionType, err := NewCompletionType(args[0])
			if err != nil {
				cli.Exit(err)
			}
			switch completionType {
			case CompletionTypeBash:
				err = rootCmd.GenBashCompletion(os.Stdout)
			case CompletionTypeZsh:
				if completionsNoDescFlag {
					err = rootCmd.GenZshCompletionNoDesc(os.Stdout)
				} else {
					err = rootCmd.GenZshCompletion(os.Stdout)
				}
			case CompletionTypeFish:
				err = rootCmd.GenFishCompletion(os.Stdout, !completionsNoDescFlag)
			case CompletionTypePowershell:
				err = rootCmd.GenPowerShellCompletion(os.Stdout)
			}
			if err != nil {
				cli.Exit(err)
			}
		},
		GroupID: "setup",
	}
	completionsCmd.Flags().BoolVar(&completionsNoDescFlag, "no-descriptions", false, "disable completions description for shells that support it")
	return &completionsCmd
}

// CompletionType defines the valid shells for which Git Town can create auto-completions.
type CompletionType string

const (
	CompletionTypeBash       CompletionType = "bash"
	CompletionTypeZsh        CompletionType = "zsh"
	CompletionTypeFish       CompletionType = "fish"
	CompletionTypePowershell CompletionType = "powershell"
)

// completionTypes provides all CompletionType values.
func completionTypes() []CompletionType {
	return []CompletionType{
		CompletionTypeBash,
		CompletionTypeZsh,
		CompletionTypeFish,
		CompletionTypePowershell,
	}
}

func NewCompletionType(text string) (CompletionType, error) {
	text = strings.ToLower(text)
	for _, completionType := range completionTypes() {
		if text == string(completionType) {
			return completionType, nil
		}
	}
	return CompletionTypeBash, fmt.Errorf("unknown completiontype: %q", text)
}
