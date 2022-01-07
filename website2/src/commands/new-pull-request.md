# New-pull-request command

```
git new-pull-request
```

The new-pull-request command allows the user to create a new pull request for
the current branch. It does that by opening a browser window showing the new
pull request page of your repository. The form is pre-populated with the current
branch and it's parent branch. This command syncs the current branch before
opening the pull request.

You can create new pull requests for repositories hosted on
[GitHub](https://github.com/), [GitLab](https://gitlab.com/),
[Gitea](https://gitea.com/) and [Bitbucket](https://bitbucket.org/). When using
self-hosted versions of these services, you can configure the hosting service
type with the [code-hosting-driver](../configurations/code-hosting-driver.md)
setting.

When using SSH identities, this command uses the hostname in the
[code-hosting-origin-hostname](../configurations/code-hosting-origin-hostname.md)
setting.
