package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionNoDesc bool

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `To load completions:

Bash:

$ source <(git-town completion bash)

# To load completions for each session, execute once:
Linux:
  $ git-town completion bash > /etc/bash_completion.d/git-town
MacOS:
  $ git-town completion bash > /usr/local/etc/bash_completion.d/git-town

Zsh:

$ source <(git-town completion zsh)

# To load completions for each session, execute once:
$ git-town completion zsh > /usr/share/zsh/vendor-completions/_git-town

Fish:

$ git-town completion fish | source

# To load completions for each session, execute once:
$ git-town completion fish > ~/.config/fish/completions/git-town.fish
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			_ = RootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			_ = RootCmd.GenZshCompletion(os.Stdout)
			// once merged https://github.com/spf13/cobra/pull/1070
			// if !completionNoDesc {
			// 	RootCmd.GenZshCompletion(os.Stdout)
			// } else {
			// 	RootCmd.GenZshCompletionNoDesc(os.Stdout)
			// }
		case "fish":
			_ = RootCmd.GenFishCompletion(os.Stdout, !completionNoDesc)
		case "powershell":
			_ = RootCmd.GenPowerShellCompletion(os.Stdout)
		}
	},
}

func init() {
	completionCmd.Flags().BoolVar(
		&completionNoDesc,
		"no-descriptions", false,
		"disable completion description for shells that support it")
	RootCmd.AddCommand(completionCmd)
}
