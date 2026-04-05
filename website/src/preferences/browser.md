# Browser

By default, Git Town launches your system's default browser by trying common
commands like `open`, `xdg-open`, or `x-www-browser`.

You can override this behavior to use a specific browser. Disable browser
launching entirely by setting `(none)` as the browser executable.

## via CLI flag

The [propose](../commands/propose.md) and [repo](../commands/repo.md) commands
allow setting the browser via the `--browser` CLI flag.

For example, to open the repo homepage using Firefox, run:

```sh
git-town repo --browser=firefox
```

## configure in config file

```toml
[hosting]
browser = "<browser executable>"
```

## configure in Git metadata

```wrap
git config [--global] git-town.browser '<browser executable>'
```

The optional `--global` flag applies this setting to all Git repositories on
your local machine. When not present, the setting applies to the current repo.

## environment variable

Git Town uses the `BROWSER` environment variable that is also used by other
tools.
