# headless

The _headless_ preference disables all interactive features.
When set, Git Town no longer:

- opens browser windows
- asks the user for data in interactive dialogs

## options

- When set to `true`, all interactive features are disabled.
- When set to `false` (default), interactive features are enabled.

## config file

In the [config file](../configuration-file.md):

```toml
[propose]
headless = true
```

## Git metadata

In Git metadata:

```bash
git config [--global] git-town.headless true
```

## environment variable

```bash
export GIT_TOWN_HEADLESS=true
```
