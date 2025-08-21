# Stash uncommitted changes

This setting configures whether Git Town commands that create a new branch
([hack](../commands/hack.md), [append](../commands/append.md),
[prepend](../commands/prepend.md) stash away uncommitted changes before creating
the new branch or not.

## options

When set to `true` (the default value), Git Town stashes away uncommitted
changes before creating the new branch and unstashes them on the new branch.
This prevents failure even in the presence of conflicts that even
`git checkout --merge` cannot resolve.

A downside of stashing uncommitted changes is that it changes what is staged and
what isn't. So if you carefully stage changes before creating new feature
branches, you can disable this option and Git Town will leave your Git index
alone.

## CLI flags

In one-off situations you can enable or disable stashing of uncommitted changes
with the `--stash` and `--no-stash` flags.

## in config file

You can disable stashing in the [config file](../configuration-file.md) like
this:

```toml
[create]
stash = false
```

## in Git metadata

To manually configure stashing in Git, run this command:

```wrap
git config [--global] git-town.stash <true|false>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, the setting applies only to the current repository.

## environment variable

You can configure whether Git Town syncs Git tags by setting the
`GIT_TOWN_STASH` environment variable.
