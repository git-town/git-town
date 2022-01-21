package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionsNoDesc bool

var completionsCmd = &cobra.Command{
	Use:   "completions [bash|zsh|fish|powershell]",
	Short: "Generates auto-completion scripts for Bash, zsh, fish, and PowerShell",
	Long: `Generates auto-completion scripts for Bash, zsh, fish, and PowerShell

With completions enabled, git-town TAB will show you all possible
subcommands. As a bonus, some shells even show the short help text next to it.

To enable completions:

Bash:

$ source <(git-town completions bash)

Persist and autoload on each session:

Linux: $ git-town completions bash > /etc/bash_completion.d/git-town

MacOS: $ git-town completions bash > /usr/local/etc/bash_completion.d/git-town

Zsh:

$ source <(git-town completions zsh)

Persist and autoload on each session:

$ git-town completions zsh > /usr/share/zsh/vendor-completions/_git-town

Fish:

CAUTION: pending upstream issue breaks this: https://github.com/spf13/cobra/pull/1122

$ git-town completions fish | source

Persist and autoload on each session:

$ git-town completions fish > /etc/fish/completions/git-town.fish

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
	RootCmd.AddCommand(completionsCmd)
}
