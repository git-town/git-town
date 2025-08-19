# GitLab token

Git Town can interact with GitLab in your name, for example to update pull
requests as branches get created, shipped, or deleted. To do so, Git Town needs
a
[personal access token](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html)
with `api` scope. You can create one in your account settings.
[account settings](https://gitlab.com/-/user_settings/personal_access_tokens) or
get one created for you by using the
[api connector type for GitLab](gitlab-connector.md).

## config file

Since your API token is confidential, you cannot add it to the config file.

## Git metadata

You can configure the API token manually by running:

```wrap
git config [--global] git-town.gitlab-token <token>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, the setting applies only to the current repository.

## environment variable

You can configure the GitLab token by setting the `GIT_TOWN_GITLAB_TOKEN`
environment variable.
