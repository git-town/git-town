# git town commit

<a type="git-town-command" />

```command-summary
git town commit [-d | --down] [--dry-run] [-h | --help] [(-m | --message) <text>] [-v | --verbose]
```

The _commit_ command commits the staged changes into another branch and syncs
these changes back into the local branch.

This helps develop changes as a stack of branches. Let's say you work on a
feature, and as part of that you discover that you need to perform some
refactoring. You want to perform the refactoring in a separate branch, so that
you can [propose](propose.md) and review it separately. You also want to build
the feature on top of the refactoring, hence the refactoring needs to happen in
an ancestor branch.

```
main
 \
  refactor
   \
    feature
```

With this branch setup, you can work on the refactor on the `feature` branch.
`git town commit` commits the changes into the `refactor` branch and syncs them
right back into the `feature` branch.

## Positional argument

When called without a positional argument, the _ship_ command ships the current
branch.

When called with a positional argument, it ships the branch with the given name.

## Options

#### `-d`<br>`--down`

When set, Git Town commits Commit into the parent branch

#### `--dry-run`

Use the `--dry-run` flag to test-drive this command. It prints the Git commands
that would be run but doesn't execute them.

#### `-h`<br>`--help`

Display help for this command.

#### `--ignore-uncommitted`<br>`--no-ignore-uncommitted`

Overrides the [ignore-uncommitted](../preferences/ignore-uncommitted.md) setting
for this call.

#### `-m <text>`<br>`--message <text>`

Similar to `git commit`, the `--message <message>` aka `-m` parameter allows
specifying the commit message via the CLI.

#### `-f <path>`<br>`--message-file <path>`

The `--message-file` aka `-f` flag uses the content of the given file for the
commit message. The filename `-` reads the commit message from STDIN.

#### `-s <name>`<br>`--strategy <name>`

Overrides the configured [ship-strategy](../preferences/ship-strategy.md).

#### `-p`<br>`--to-parent`

The _ship_ command ships only direct children of the main branch. To ship a
child branch, you need to first ship or [delete](delete.md) all its ancestor
branches. If you really want to ship into a non-perennial branch, you can
override the protection against that with the `--to-parent` aka `-p` option.

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.

## Configuration

The configured [ship-strategy](../preferences/ship-strategy.md) determines how
the _ship_ command merges branches. When shipping
[stacked changes](../stacked-changes.md), use the
[fast-forward ship strategy](../preferences/ship-strategy.md#fast-forward) to
avoid empty merge conflicts.

If you have configured the API tokens for
[GitHub](../preferences/github-token.md),
[GitLab](../preferences/gitlab-token.md),
[Gitea](../preferences/gitea-token.md),
[Bitbucket](../preferences/bitbucket-app-password.md), or
[Forgejo](../preferences/forgejo-token.md) and the branch to be shipped has an
open proposal, this command merges the proposal for the current branch.

If your forge automatically deletes shipped branches, for example
[GitHub's feature to automatically delete head branches](https://help.github.com/en/github/administering-a-repository/managing-the-automatic-deletion-of-branches),
you can
[disable deleting remote branches](../preferences/ship-delete-tracking-branch.md).

## See also

<!-- keep-sorted start -->

- [propose](propose.md) creates a pull request for the current branch
- [repo](repo.md) opens the website of your forge in the browser, so that you
  can ship branches there

<!-- keep-sorted end -->
