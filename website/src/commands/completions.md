# git town completions

```command-summary
git town completion <bash|fish|powershell|zsh> [--no-descriptions]
```

The _completions_ command outputs shell scripts that enable auto-completion for
Git Town in Bash, Zsh, Fish, or PowerShell. When set up, typing
`git-town <tab key>` in your terminal will auto-complete subcommands.

## Options

#### `--no-descriptions`

The `--no-descriptions` flag outputs shorter completions without descriptions of
arguments.

### Bash

To load autocompletion for Bash, run this command:

```
source <(git-town completions bash)
```

To load completions for each session, add the above line to your `.bashrc`.

### Fish

To load autocompletions for Fish, run this command:

```
git-town completions fish | source
```

To load completions for each session, add the above line to your
`~/.config/fish/config.fish`.

### PowerShell

To install autocompletions for Powershell, run this command:

```
git-town completions powershell | Out-String | Invoke-Expression
```

To load completions for each session, add the above line to your PowerShell
profile.

### Zsh

To load autocompletions for Zsh, run this command:

```
source <(git-town completions zsh)
```

To load completions for each session, add the above line to your `.zshrc`.

To fix the error message `command not found: compdef`, run

```zsh
autoload -Uz compinit
```
