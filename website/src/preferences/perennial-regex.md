# pererennial-regex

All branches matching this regular expression are considered
[perennial branches](perennial-branches.md).

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
git config [--global] git-town.perennial-regex 'release-.*'
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, the setting applies only to the current repository.
