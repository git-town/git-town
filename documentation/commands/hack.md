<h1 textrun="command-heading">Hack command</h1>

<blockquote textrun="command-summary">
Creates a new feature branch off the main development branch
</blockquote>

<a textrun="command-description">
Syncs the main branch,
forks a new feature branch with the given name off the main branch,
pushes the new feature branch to the remote repository
(if and only if [new-branch-push-flag](./new-branch-push-flag.md) is true),
and brings over all uncommitted changes to the new feature branch.

See [sync](./sync.md) for information regarding remote upstream. </a>

#### Usage

<pre textrun="command-usage">
git town hack &lt;branch&gt;
</pre>
