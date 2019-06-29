<h1 textrun="command-heading">Hack command</h1>

<blockquote textrun="command-summary">
Creates a new feature branch off the main development branch
</blockquote>

<a textrun="command-description">
Syncs the main branch and forks a new feature branch with the given name off it.

If (and only if) [new-branch-push-flag](./new-branch-push-flag.md) is true,
pushes the new feature branch to the remote repository.

Finally, brings over all uncommitted changes to the new feature branch.

Additionally, when there is a remote upstream,
the main branch is synced with its upstream counterpart.
This can be disabled by toggling the "new-branch-push-flag" configuration.
</a>

#### Usage

<pre textrun="command-usage">
git town hack &lt;branch&gt;
</pre>
