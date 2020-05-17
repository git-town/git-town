<h1 textrun="command-heading">Repo command</h1>

<blockquote textrun="command-summary">
Opens the repository homepage
</blockquote>

<a textrun="command-description">

Supported for repositories hosted on [GitHub](http://github.com/),
[GitLab](http://gitlab.com/), and [Bitbucket](https://bitbucket.org/). Derives
the Git provider from the `origin` remote. You can override this detection with
`git config git-town.code-hosting-driver <DRIVER>` where DRIVER is "github",
"gitlab", or "bitbucket".

When using SSH identities, run
`git config git-town.code-hosting-origin-hostname <HOSTNAME>` where HOSTNAME
matches what is in your ssh config file.

</a>

#### Usage

<pre textrun="command-usage">
git town repo
</pre>
