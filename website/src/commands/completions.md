# git town completions

> _git town completion <bash|zsh|fish|powershell>_

The _completions_ command outputs shell scripts that enable auto-completion for
Git Town in Bash, Zsh, Fish, or PowerShell. When set up, typing
`git-town <tab key>` in your terminal will auto-complete subcommands.

## bash

To load autocompletion for Bash, run this command:

```
source <(git-town completions bash)
```

To load completions for each session, add the above line to your `.bashrc`.

## zsh

To load autocompletions for Zsh, run this command:

```
source <(git-town completions zsh)
```

To load completions for each session, add the above line to your `.zshrc`.

To fix the error message `command not found: compdef`, run

```zsh
autoload -Uz compinit
```

## fish

To load autocompletions for Fish, run this command:

```
git-town completions fish | source
```

To load completions for each session, add the above line to your
`~/.config/fish/config.fish`.

## powershell

To install autocompletions for Powershell, run this command:

```
git-town completions powershell | Out-String | Invoke-Expression
```

To load completions for each session, add the above line to your PowerShell
profile.
