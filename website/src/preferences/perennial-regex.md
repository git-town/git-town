# pererennial-regex

If you have many [perennial branches](perennial-branches.md) with comparable
names, you can define them via a regular expression.

## configure in config file

In the [config file](../configuration-file.md) the perennial regex exists inside
the `[branches]` section:

```toml
[branches]
perennial-regex = "release-*"
```

## configure in Git metadata

You can configure the perennial branches manually by running:

```bash
git config [--global] git-town.perennial-branches "branch other-branch"
```

The optional `--global` flag applies this setting to all Git repositories on
your local machine. When not present, the setting applies to the current repo.
