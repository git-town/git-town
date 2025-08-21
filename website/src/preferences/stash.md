# Stash uncommitted changes

This setting controls whether Git Town stashes uncommitted changes away before
creating and switching to a new branch, i.e. affects the
[hack](../commands/hack.md), [append](../commands/append.md), and
[prepend](../commands/prepend.md) commands.

## options

By default (`true`), Git Town stashes your uncommitted changes before creating
the new branch and restores them afterwards. This ensures the branch switch
succeeds, even if there are conflicts that `git checkout --merge` can't handle.

The tradeoff is that if you had changes stashed before, those changes are now
unstashed. If you carefully staged changes before creating a new branch, you may
want to disable this option to keep your index untouched.

## CLI flags

You can override this setting per command using:

- `--stash` to force stashing
- `--no-stash` to skip stashing

## in config file

To permanently disable stashing in the [config file](../configuration-file.md):

```toml
[create]
stash = false
```

## in Git metadata

You can also configure stashing only on your machine:

```wrap
git config [--global] git-town.stash <true|false>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, the setting applies only to the current repository.

## environment variable

The `GIT_TOWN_STASH` environment variable also configures this behavior.
