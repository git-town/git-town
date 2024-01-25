# Quick configuration

Git Town prompts for all the configuration information it needs. The commands
below set additional configuration options that might be helpful in your use
case. We will cover the full list of options [later](configuration-commands.md).

## API access to your hosting provider

Git Town can ship branches that have an open proposal by merging this proposal
via your code hosting service's API. To enable it, create an API token for your
account at your code hosting platform.

- [instructions for GitHub](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token)
- [instructions for GitLab](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html)
- [instructions for Gitea](https://docs.gitea.io/en-us/api-usage)

Provide the token to Git Town as an environment variable with name
`GITHUB_TOKEN` or `GITHUB_AUTH_TOKEN` or as part of the Git Town configuration:

- GitHub:

  ```
  git config --add git-town.github-token <your api token>
  ```

- GitLab:

  ```
  git config --add git-town.gitlab-token <your api token>
  ```

## Delete tracking branches

Some code hosting platforms
[automatically delete feature branches](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/configuring-pull-request-merges/managing-the-automatic-deletion-of-branches)
after merging them. This is a useful feature that you should enable if possible.
It can interfere with Git Town's attempts to also delete this branch after
shipping it. To make Git Town play along, run:

```
git config git-town.ship-delete-tracking-branch false
```

## Shell autocompletion

Follow the instructions given by `git-town help completions` to install the
autocompletions for your shell.

## Shorter commands

Having to type `git town <command>` gets old. Git Town can install aliases for
its commands that make them feel like native Git commands, i.e. allow you to run
for example `git hack` instead of `git town hack`. To enable this feature:

```
git town aliases add
```

To remove these aliases, run `git town aliases remove`.
