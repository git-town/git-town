<h1 textrun="command-heading">Prepend command</h1>

<blockquote textrun="command-summary">
Creates a new feature branch as the parent of the current branch
</blockquote>

<a textrun="command-description">

Syncs the parent branch, cuts a new feature branch with the given name off the
parent branch, makes the new branch the parent of the current branch, pushes the
new feature branch to the remote repository (if
[new-branch-push-flag](./new-branch-push-flag.md) is true), and brings over all
uncommitted changes to the new feature branch.

See [sync](./sync.md) for remote upstream options.

</a>

#### Usage

<pre textrun="command-usage">
git town prepend &lt;branch&gt;
</pre>

#### SEE ALSO

- [git append](append.md) to create a new feature branch as a child of the
  current branch
- [git hack](hack.md) to create a new top-level feature branch
