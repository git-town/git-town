package cmd

import (
	"os"

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
			switch args[0] {
			case "bash":
				_ = rootCmd.GenBashCompletion(os.Stdout)
			case "zsh":
				if completionsNoDescFlag {
					_ = rootCmd.GenZshCompletionNoDesc(os.Stdout)
				} else {
					_ = rootCmd.GenZshCompletion(os.Stdout)
				}
			case "fish":
				_ = rootCmd.GenFishCompletion(os.Stdout, !completionsNoDescFlag)
			case "powershell":
				_ = rootCmd.GenPowerShellCompletion(os.Stdout)
			}
		},
		GroupID: "setup",
	}
	completionsCmd.Flags().BoolVar(&completionsNoDescFlag, "no-descriptions", false, "disable completions description for shells that support it")
	return &completionsCmd
}
