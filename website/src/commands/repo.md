# git town repo

> _git town repo [remote-name]_

The _repo_ command ("show me the repository") opens the homepage of the current
repository in your browser. Git Town can display repositories hosted on
[GitHub](https://github.com), [GitLab](https://gitlab.com),
[Gitea](https://gitea.com), and [Bitbucket](https://bitbucket.org).

### Positional arguments

When called without arguments, the _repo_ command shows the repository at the
[development remote](../preferences/dev-remote.md).

When called with an argument, it shows the repository at the remote with the
given name.

### Configuration

Git Town automatically identifies the hosting platform type through the URL of
the development remote. You can override the type of hosting server with the
[hosting-platform](../preferences/hosting-platform.md) setting.

Set the [hosting-origin-hostname](../preferences/hosting-origin-hostname.md)
setting to tell Git Town about the hostname when using ssh identities.

### --verbose / -v

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.
