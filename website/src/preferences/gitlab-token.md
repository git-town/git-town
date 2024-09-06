# gitlab-token

Git Town can interact with GitLab in your name, for example to update pull
requests as branches get created, shipped, or deleted. To do so, Git Town needs
a
[personal access token](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html)
with `api` scope. You can create one in your account settings.

The best way to enter your token is via the
[setup assistant](../configuration.md).

## config file

Since your API token is confidential, you cannot add it to the config file.

## Git metadata

You can configure the API token manually by running:

```bash
git config [--global] git-town.gitlab-token <token>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, the setting applies only to the current repository.
