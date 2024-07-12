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

GitHub supports these additional command-line switches:

- `--title` pre-populate the title of the pull request
- `--body` pre-populates the body of the pull request
- `--body-file` pre-populates the body of the pull request with the content of
  the given file. The filename `-` makes Git Town read the body text from STDIN.

### Configuration

You can configure the hosting platform type with the
[hosting-platform](../preferences/hosting-platform.md) setting.

When using SSH identities, this command uses the hostname in the
[hosting-origin-hostname](../preferences/hosting-origin-hostname.md) setting.
