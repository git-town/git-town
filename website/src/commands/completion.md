# git town completion [bash|zsh|fish|powershell]

The _completion_ command generates auto-completion scripts for Bash, Zsh, Fish,
and PowerShell. With shell completions set up, typing `git-town <tab key>` in
your terminal will auto-complete subcommands.

## Bash

To install autocompletion for Bash, run this command:

```
source <(git-town completion bash)
```

## Zsh

To install autocompletions for Zsh, run this command:

```
source <(git-town completion zsh)
```

## Fish

To install autocompletions for Fish, run this command:

```
git-town completion fish | source
```

## Powershell

To install autocompletions for Powershell, run this command:

```
Invoke-Expression -Command $(git-town completion powershell | Out-String)
```
