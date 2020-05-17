<h1 textrun="command-heading">Ship command</h1>

<blockquote textrun="command-summary">
Deliver a completed feature branch
</blockquote>

<a textrun="command-description">

Squash-merges the current branch, or <branch_name> if given, into the main
branch, resulting in linear history on the main branch.

- syncs the main branch
- pulls remote updates for <branch_name>
- merges the main branch into <branch_name>
- squash-merges <branch_name> into the main branch with commit message specified
  by the user
- pushes the main branch to the remote repository
- deletes <branch_name> from the local and remote repositories

Ships direct children of the main branch. To ship a nested child branch, ship or
kill all ancestor branches first.

If you use GitHub, this command can squash merge pull requests via the GitHub
API. Setup:

1. Get a
   [GitHub personal access token](https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line)
   with the "repo" scope
2. Run 'git config git-town.github-token XXX' (optionally add the '--global'
   flag) Now anytime you ship a branch with a pull request on GitHub, it will
   squash merge via the GitHub API. It will also update the base branch for any
   pull requests against that branch.

If your origin server deletes shipped branches, for example
[GitHub's feature to automatically delete head branches](https://help.github.com/en/github/administering-a-repository/managing-the-automatic-deletion-of-branches),
run `git config git-town.ship-delete-remote-branch false` and Git Town will
leave it up to your origin server to delete the remote branch.

</a>

#### Usage

<pre textrun="command-usage">
git town ship
</pre>

#### Flags

<pre textrun="command-flags">
-m, --message string   Specify the commit message for the squash commit
</pre>
