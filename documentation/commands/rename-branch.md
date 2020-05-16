<h1 textrun="command-heading">Rename-branch command</h1>

<blockquote textrun="command-summary">
Renames a branch both locally and remotely
</blockquote>

<a textrun="command-description">
Renames the given branch in the local and origin repository.
Aborts if the new branch name already exists or the tracking branch is out of sync.

- creates a branch with the new name
- deletes the old branch

When there is a remote repository

- syncs the repository

When there is a tracking branch

- pushes the new branch to the remote repository
- deletes the old branch from the remote repository

When run on a perennial branch

- confirm with the "-f" option
- registers the new perennial branch name in the local Git Town configuration

</a>

#### Usage

<pre textrun="command-usage">
git town rename-branch [&lt;old_branch_name&gt;] &lt;new_branch_name&gt;
</pre>

#### Flags

<pre textrun="command-flags">
--force   Force rename of perennial branch
</pre>
