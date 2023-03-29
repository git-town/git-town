# git repo

The _repo_ command ("show the repository") opens the homepage of the current
repository in your default browser. Git Town can display repositories hosted on
[GitHub](https://github.com), [GitLab](https://gitlab.com),
[Gitea](https://gitea.com), [Bitbucket](https://bitbucket.org), and
[Azure DevOps](https://azure.microsoft.com/en-us/products/devops).

### Variations

Git Town identifies the hosting service type by looking at the `origin` remote.
You can override this detection with the
[code-hosting-driver](../preferences/code-hosting-driver.md) setting.

Set the
[code-hosting-origin-hostname](../preferences/code-hosting-origin-hostname.md)
setting to tell Git Town about the hostname when using ssh identities.
