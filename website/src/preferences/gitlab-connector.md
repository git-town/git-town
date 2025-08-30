# GitLab connector

Git Town can interact with GitLab in two different ways.

1. **GitLab API:** <br> Git Town communicates directly with the GitLab API using
   a personal access token. You'll need to generate this token at
   [gitlab.com/settings/tokens](https://gitlab.com/-/user_settings/personal_access_tokens).

   By default, Git stores such tokens in plaintext in your Git configuration. To
   avoid this, consider configuring Git to use your operating system's encrypted
   [credentials storage](https://git-scm.com/book/en/v2/Git-Tools-Credential-Storage)
   for better security.

2. **GitLab CLI (glab):** <br> The [glab](https://gitlab.com/gitlab-org/cli) CLI
   handles authentication and token management for you.

## config file

It is generally not recommended to hardcode the connector type in your config
file, as it enforces usage or non-usage of `glab` for your entire team. If you
want to set it explicitly, it would look like this:

```toml
[hosting]
gitlab-connector = "api" # or "glab"
```

## Git metadata

You can configure the API token manually by running:

```wrap
git config [--global] git-town.gitlab-connector <api|glab>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, the setting applies only to the current repository.

## environment variable

You can configure the GitLab connector by setting the
`GIT_TOWN_GITLAB_CONNECTOR` environment variable.
