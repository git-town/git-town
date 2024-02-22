# observed-branches

Observed branches are branches that you have checked out on your machine, and
might commit into, but you don't want to push your commits to origin. This is
useful for staff-level engineers and tech leads that support other teams

You can see the configured perennial branches via the
[config](../commands/config.md) command and change them via the
[setup assistant](../commands/config-setup.md).

## configure in config file

In the [config file](../configuration-file.md) the perennial branches are
defined as part of the `[branches]` section:

```toml
[branches]
perennials = [ "branch", "other-branch" ]
```

## configure in Git metadata

You can configure the perennial branches manually by running:

```bash
git config [--global] git-town.perennial-branches "branch other-branch"
```

The optional `--global` flag applies this setting to all Git repositories on
your local machine. When not present, the setting applies to the current repo.

## bulk-define perennial branches

If you have many perennial branches that follow the same naming schema, like
`release-v4.0-rev.1`, `release-v4.0-rev.2`, etc, you can define a
[regular expression](perennial-regex.md) for them instead of listing them one by
one.
