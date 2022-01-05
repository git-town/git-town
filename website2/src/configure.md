# Configuration

Here are a few options to configure Git Town after you have
[installed](install.md) it.

### Enable shorter commands

Having to type `git town <command>` gets old. Git Town can install aliases for
its commands that make them feel like native Git commands, i.e. you can run
`git hack` instead of `git town hack` etc. To enable this feature:

```
git town alias true
```

To remove the aliases, run `git town alias false`.

### Enable API access to your hosting provider

If you host on GitHub or GitTea, Git Town can merge pull requests via the API of
your code hosting service. This is better than merging feature branches locally
because it merges rather than closes the affected pull request. To enable this
feature, create an API token for your account at your code hosting provider:

- [instructions for GitHub](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token)
- [instructions for Gitea](https://docs.gitea.io/en-us/api-usage)

Once you have an API token for your account, tell Git Town about it by adding it
to your Git configuration:

```
git config --add git-town.github-token <your api token>
```

Once set up, Git Town ships feature branches that have an open pull request via
the API of the hosting service.

### Enable automatically deleted head branches

Some code hosting providers
[automatically delete feature branches](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/configuring-pull-request-merges/managing-the-automatic-deletion-of-branches)
after merging them. This is a very useful feature that you should enable if
possible, but it can trip up Git Town when it also tries to delete this branch
after shipping it. To make Git Town play along, run:

```
git config git-town.ship-delete-remote-branch false
```

### Install autocompletion

Follow the instructions given by `git-town help completions` to install the
autocompletions for your shell.

```
```
