<h1 textrun="command-heading">Prepend command</h1>

<blockquote textrun="command-summary">
Creates a new feature branch as the parent of the current branch
</blockquote>

<a textrun="command-description">
Syncs the parent branch,
forks a new feature branch with the given name off the parent branch,
makes the new branch the parent of the current branch,
pushes the new feature branch to the remote repository,
and brings over all uncommitted changes to the new feature branch.

Additionally, when there is a remote upstream,
the main branch is synced with its upstream counterpart.
This can be disabled by toggling the "new-branch-push-flag" configuration:

```
git town new-branch-push-flag false
```

</a>

#### Usage

<pre textrun="command-usage">
git town prepend &lt;branch&gt;
</pre>

#### SEE ALSO

- [git append](append.md) to create a new feature branch as a child of the current branch
- [git hack](hack.md) to create a new top-level feature branch
