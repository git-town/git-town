package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionsNoDesc bool

var completionsCmd = &cobra.Command{
	Use:   "completions [bash|zsh|fish|powershell]",
	Short: "Generates auto-completion for bash, zsh, fish, or PowerShell",
	Long: `Generates auto-completion for bash, zsh, fish, or PowerShell

When set up, "git-town <TAB>" will auto-complete Git Town subcommands.

To enable completions:

Bash:

$ source <(git-town install completions bash)

Persist and autoload on each session:

Linux: $ git-town install completions bash > /etc/bash_completion.d/git-town

MacOS: $ git-town install completions bash > /usr/local/etc/bash_completion.d/git-town

Zsh:

$ source <(git-town install completions zsh)

Persist and autoload on each session:

$ git-town install completions zsh > /usr/share/zsh/vendor-completions/_git-town

Fish:

$ git-town install completions fish | source

Persist and autoload on each session:

$ git-town install completions fish > /etc/fish/completions/git-town.fish

You might be a power user who has their dotfiles under version control. Or you
might have another motivation to keep those scripts in your home folder.
Since it sometimes depends on your particular setup, you probably should consult
the official docs for your shell.
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
