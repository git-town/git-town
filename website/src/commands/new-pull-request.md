# git new-pull-request

The _new-pull-request_ command helps create a new pull request for the current
feature branch. It does that by opening a browser window showing the new pull
request page of your code hosting service. The form is prepopulated with the
current branch and it's parent branch. This command [syncs](sync.md) the current
branch before opening the pull request.

### Variations

You can create new pull requests for repositories hosted on:

- [Bitbucket](https://bitbucket.org)
- [Gitea](https://gitea.com)
- [GitHub](https://github.com)
- [GitLab](https://gitlab.com)

When using self-hosted versions of these services, you can configure the hosting
service type with the
[code-hosting-driver](../preferences/code-hosting-driver.md) setting.

When using SSH identities, this command uses the hostname in the
[code-hosting-origin-hostname](../preferences/code-hosting-origin-hostname.md)
setting.
