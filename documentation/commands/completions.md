<h1 textrun="command-heading">Completions command</h1>

<blockquote textrun="command-summary">
Generates completion scripts for Bash, zsh, fish, and PowerShell
</blockquote>

<a textrun="command-description">
A Git Town user values productivity, so this is for you.

With completions enabled, git-town TAB will show you all possible subcommands.
As a bonus, some shells even show the short help text next to it.

To enable completions:

Bash:

\$ source <(git-town completions bash)

Persist and autoload on each session:

Linux: \$ git-town completions bash > /etc/bash_completion.d/git-town

MacOS: \$ git-town completions bash > /usr/local/etc/bash_completion.d/git-town

Zsh:

\$ source <(git-town completions zsh)

Persist and autoload on each session:

\$ git-town completions zsh > /usr/share/zsh/vendor-completions/\_git-town

Fish:

\$ git-town completions fish | source

Persist and autoload on each session:

\$ git-town completions fish > /etc/fish/completions/git-town.fish

You might be a power user who has their dotfiles under version control. Or you
might have another motivation to keep those scripts in your home folder. Since
it sometimes depends on your particular setup, you probably should consult the
official docs for your shell.</a>

#### Usage

<pre textrun="command-usage">
git town completions <bash|zsh|fish|powershell>
</pre>
