<h1 textrun="command-heading">New-pull-request command</h1>

<blockquote textrun="command-summary">
Creates a new pull request
</blockquote>

<a textrun="command-description">
Syncs the current branch
and opens a browser window to the new pull request page of your repository.
The form is pre-populated for the current branch
so that the pull request only shows the changes made
against the immediate parent branch.

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
git town new-pull-request
</pre>
