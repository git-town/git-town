package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionsNoDesc bool

var completionsCmd = &cobra.Command{
	Use:   "completions [bash|zsh|fish|powershell]",
	Short: "Generates auto-completion for bash, zsh, fish, or PowerShell",
	Long: `Generates auto-completion for bash, zsh, fish, or PowerShell.
When set up, "git-town <TAB>" will auto-complete Git Town subcommands.

To load autocompletion for Bash, run this command:

	git-town install completions bash | source

To load completions for each session, add the above line to your ~/.bashrc file.


To load autocompletion for Zsh, run this command:

	source <(git-town install completions zsh)

To load completions for each session, add the above line to your ~/.zshrc file.


To load autocompletion for Fish, run this command:

	git-town install completions fish | source

To load completions for each session, add the above line to your ~/.config/fish/config.fish file.


To load autocompletions for Powershell, run this command:

	git-town install completions powershell | Out-String | Invoke-Expression

To load completions for each session, add the above line to your PowerShell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			_ = RootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			if completionsNoDesc {
				_ = RootCmd.GenZshCompletionNoDesc(os.Stdout)
			} else {
				_ = RootCmd.GenZshCompletion(os.Stdout)
			}
		case "fish":
			_ = RootCmd.GenFishCompletion(os.Stdout, !completionsNoDesc)
		case "powershell":
			_ = RootCmd.GenPowerShellCompletion(os.Stdout)
		}
	},
}

func init() {
	completionsCmd.Flags().BoolVar(
		&completionsNoDesc,
		"no-descriptions", false,
		"disable completions description for shells that support it")
	installCommand.AddCommand(completionsCmd)
	RootCmd.CompletionOptions.DisableDefaultCmd = true
}
