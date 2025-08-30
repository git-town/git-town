# Main branch

This setting stores the name of the main branch. The main branch is the default
parent branch for new feature branches created with
[git town hack](../commands/hack.md) and the default branch into which Git Town
[ships](../commands/ship.md) finished feature branches.

## config file

In the [config file](../configuration-file.md) the main branch is part of the
`[branches]` section:

```toml
[branches]
main = "config-main"
```

## Git metadata

To configure the main branch in Git, run this command:

```wrap
git config [--global] git-town.main-branch <value>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, the setting applies only to the current repository.

## environment variable

You can configure the main branch by setting the `GIT_TOWN_MAIN_BRANCH`
environment variable.
