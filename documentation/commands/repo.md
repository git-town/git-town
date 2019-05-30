<h1 textrun="command-heading">Repo command</h1>

<blockquote textrun="command-summary">
Opens the repository homepage
</blockquote>

<a textrun="command-description">

Supported only for repositories hosted on [GitHub](http://github.com/),
[GitLab](http://gitlab.com/), and [Bitbucket](https://bitbucket.org/).
When using self-hosted versions this command needs to be configured with
`git config git-town.code-hosting-driver <driver>`
where driver is "github", "gitlab", or "bitbucket".
When using SSH identities, this command needs to be configured with
`git config git-town.code-hosting-origin-hostname <hostname>`
where hostname matches what is in your ssh config file.
</a>

#### Usage

<pre textrun="command-usage">
git town repo
</pre>
