# Perennial regex

All branches matching this regular expression are considered
[perennial branches](perennial-branches.md).

## configure in config file

In the [config file](../configuration-file.md) the perennial regex exists inside
the `[branches]` section:

```toml
[branches]
perennial-regex = "^release-.*"
```

## configure in Git metadata

You can configure the perennial branches manually by running:

```wrap
git config [--global] git-town.perennial-regex 'release-.*'
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, the setting applies only to the current repository.

## environment variable

You can configure the perennial regex by setting the `GIT_TOWN_PERENNIAL_REGEX`
environment variable.
