<h1 textrun="command-heading">Append command</h1>

<blockquote textrun="command-summary">
Creates a new feature branch as a direct child of the current branch.
</blockquote>

<a textrun="command-description">
Syncs the current branch,
forks a new feature branch with the given name off the current branch,
makes the new branch a child of the current branch,
pushes the new feature branch to the remote repository
if and only if <a href="./new-branch-push-flag.md">new-branch-push-flag</a>
is true,
and brings over all uncommitted changes to the new feature branch.
</a>

#### Usage

<pre textrun="command-usage">
git town append &lt;branch&gt;
</pre>
