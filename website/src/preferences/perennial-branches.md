# pererennial-branches

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
