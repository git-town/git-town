# git town repo

```command-summary
git town repo [<remote-name>] [-v | --verbose]
```

The _repo_ command ("show me the repository") opens the homepage of the current
repository in your browser. Git Town can display repositories hosted on
[GitHub](https://github.com), [GitLab](https://gitlab.com),
[Gitea](https://gitea.com), [Bitbucket](https://bitbucket.org), and
[Codeberg](https://codeberg.org).

On non-Windows systems, Git Town will first read the `BROWSER` environment
variable to determine the browser command. If it isn't set, Git Town will try
various common commands like `open`, `xdg-open`, or `x-www-browser`.

## Positional arguments

When called without arguments, the _repo_ command shows the repository at the
[development remote](../preferences/dev-remote.md).

When called with an argument, it shows the repository at the remote with the
given name.

## Options

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.

## Configuration

Git Town automatically identifies the forge type through the URL of the
development remote. You can override the type of hosting server with the
[hosting-platform](../preferences/forge-type.md) setting.

Set the [hosting-origin-hostname](../preferences/hosting-origin-hostname.md)
setting to tell Git Town about the hostname when using ssh identities.
