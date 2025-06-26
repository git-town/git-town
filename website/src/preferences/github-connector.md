# GitHub connector

Git Town can interact with GitHub in two different ways.

1. **GitHub API:** Git Town talks directly with the GitHub API. This uses an
   access token for your account, which you need to create at
   [github.com/settings/tokens](https://github.com/settings/tokens). By default,
   Git stores this token in clear text in your Git configuration. You can
   connect Git to your operating system's encrypted
   [credentials storage](https://git-scm.com/book/en/v2/Git-Tools-Credential-Storage)
   for better security.

2. **GitHub's "gh" application:** This sets up and stores the access token for
   you, but you need to install and configure the [gh](https://cli.github.com)
   tool.

## config file

Storing this value in the config file is not recommended because it forces all
members of your team to use this connector type. Having said that, the config
file would contain this value like this:

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
