# Set up configuration

If your repository already contains a `.git-branches.toml` or `.git-town.toml`
file, you are good to go. If not, or something doesn't work, you can run Git
Town's setup assistant to configure Git Town. It walks you through every
available configuration option, explains it, and gives you a chance to adjust
it.

```
git town config setup
```

You can find more background around how Git Town stores configuration in the
[overview of all configuration options](preferences.md).

### Access tokens

Some of Git Towns' functionality requires access the API of your forge:

- if the parent of a branch is not known, Git Town can look for a pull requests
  of this branch and uses their parent branch
- when you prepend, rename, remove branches or change their parent, Git Town can
  updates the affected pull requests
- click the "merge" button on a pull request from your CLI

Configuring API access is easy. Here is how you do it:

- GitHub: You can set up the [access token](preferences/github-token.md) that
  allows Git Town to talk to GitHub's API, or set up GitHub's
  [gh tool](https://cli.github.com) and let Git Town interact with Github
  through it
- GitLab: [access token](preferences/gitlab-token.md)
- Bitbucket: [username](preferences/bitbucket-username.md) and
  [app password](preferences/bitbucket-app-password.md)
- Gitea: [access token](preferences/gitea-token.md)
- Codeberg: [access token](preferences/codeberg-token.md)
