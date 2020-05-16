<h1 textrun="command-heading">Sync command</h1>

<blockquote textrun="command-summary">
Updates the current branch with all relevant changes
</blockquote>

<a textrun="command-description">
Synchronizes the current branch with the rest of the world.

When run on a feature branch:

- syncs all ancestor branches
- pulls updates for the current branch
- merges the parent branch into the current branch
- pushes the current branch

When run on the main branch or a perennial branch:

- pulls and pushes updates for the current branch
- pushes tags

If the repository contains an "upstream" remote, syncs the main branch with its
upstream counterpart. You can disable this by running
`git config git-town.sync-upstream false`.

</a>

#### Usage

<pre textrun="command-usage">
git town sync
</pre>

#### Flags

<pre textrun="command-flags">
--all       Sync all local branches
--dry-run   Print the commands but don't run them
</pre>
