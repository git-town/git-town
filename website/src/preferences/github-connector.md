# GitHub connector

Git Town can interact with GitHub in two different ways.

1. **GitHub API:** <br> Git Town communicates directly with the GitHub API using
   a personal access token. You'll need to generate this token
   [github.com/settings/tokens](https://github.com/settings/tokens).

   By default, Git stores such tokens in plaintext in your Git configuration. To
   avoid this, consider configuring Git to use your operating system's encrypted
   [credentials storage](https://git-scm.com/book/en/v2/Git-Tools-Credential-Storage)
   for better security.

2. **GitHub CLI (gh):** <br> The [gh](https://cli.github.com) CLI handles
   authentication and token management for you.

## config file

It is generally not recommended to hardcode the connector type in your config
file, as it enforces usage or non-usage of `gh` for your entire team. If you
want to set it explicitly, it would look like this:

```toml
[hosting]
github-connector = "api" # or "gh"
```

## Git metadata

You can configure the API token manually by running:

```wrap
git config [--global] git-town.github-connector <api|gh>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, the setting applies only to the current repository.

## environment variable

You can configure the GithubConnector by setting the `GIT_TOWN_GITHUB_CONNECTOR`
environment variable.
