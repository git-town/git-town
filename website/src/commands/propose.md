# git propose

The _propose_ command helps create a new pull request (also known as merge
request) for the current feature branch. It opens your code hosting platform's
website to create a new proposal in your browser and pre-populates information
like branch and source/target repository. It also [syncs](sync.md) the branch to
merge before opening the pull request.

You can create pull requests for repositories hosted on:

- [Bitbucket](https://bitbucket.org)
- [Gitea](https://gitea.com)
- [GitHub](https://github.com)
- [GitLab](https://gitlab.com)

### Arguments

On GitHub you can pre-populate the title of the pull request with the `--title`
argument and the body text of the pull request with the `--body` argument.

### Configuration

You can configure the hosting platform type with the
[hosting-platform](../preferences/hosting-platform.md) setting.

When using SSH identities, this command uses the hostname in the
[hosting-origin-hostname](../preferences/hosting-origin-hostname.md) setting.
