<h1 textrun="command-heading">Completions command</h1>

<blockquote textrun="command-summary">
Generates completion scripts for Bash, zsh, fish, and PowerShell
</blockquote>

<a textrun="command-description">
Shell completions are the productuvity boost you presumably are
after as a Git Town user.

With completions enabled, `git-town <TAB>` will show you all possible
subcommands. As a bonus, some shells even the short help text next to it.

## To enable completions:

### Bash:

`$ source <(git-town completions bash)`

**Persist and autoload on each session:**

Linux: `$ git-town completions bash > /etc/bash_completion.d/git-town`

MacOS: `$ git-town completions bash > /usr/local/etc/bash_completion.d/git-town`

### Zsh:

`$ source <(git-town completions zsh)`

**Persist and autoload on each session:**
`$ git-town completions zsh > /usr/share/zsh/vendor-completions/_git-town`

### Fish:

`$ git-town completions fish | source`

**Persist and autoload on each session:**
`$ git-town completions fish > ~/.config/fish/completions/git-town.fish` </a>

#### Usage

<pre textrun="command-usage">
git town completions [bash|zsh|fish|powershell]
</pre>
