# perennial-branches

Perennial branches are long-lived branches. They have no parent and are never
shipped. Typical perennial branches are `main`, `master`, `development`,
`production`, `staging`, etc.

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
your machine. Without it, the setting applies only to the current repository.

## bulk-define perennial branches

If you have many perennial branches that follow the same naming schema, like
`release-v4.0-rev.1`, `release-v4.0-rev.2`, etc, you can define a
[regular expression](perennial-regex.md) for them instead of listing them one by
one.
