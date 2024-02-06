# gitlab-token

To interact with the GitLab API, Git Town needs a
[personal access token](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html)
with `api` scope.

The best way to enter your token is via the
[setup assistant](../configuration.md). Since your API token is confidential,
you cannot enter it into the config file.

## in Git metadata

You can configure the API token manually by running:

```bash
git config [--global] git-town.gitlab-token <token>
```
