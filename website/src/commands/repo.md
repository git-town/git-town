# git repo [remote]

The _repo_ command ("show the repository") opens the homepage of the current
repository in your default browser. Git Town can display repositories hosted on
[GitHub](https://github.com), [GitLab](https://gitlab.com),
[Gitea](https://gitea.com), and [Bitbucket](https://bitbucket.org).

### Arguments

Shows the repository at the remote with the given name. If no remote is given,
shows the repository at the `origin` remote.

### Configuration

Git Town automatically identifies the hosting platform type through the `origin`
remote. You can override the type of hosting server with the
[hosting-platform](../preferences/hosting-platform.md) setting.

Set the [hosting-origin-hostname](../preferences/hosting-origin-hostname.md)
setting to tell Git Town about the hostname when using ssh identities.
