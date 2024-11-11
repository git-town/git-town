# git town propose

> _git town propose [--title &lt;text&gt;] [--body &lt;text&gt;] [--body-file
> &lt;-|filename&gt;]_

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

### --body / -b

When called with the `--body` aka `-b` flag, it pre-populates the body of the
pull request with the given text.

### --body-file / -f

When called with the `--body-file` aka `-f` flag, it pre-populates the body of
the pull request with the content of the given file. The filename `-` reads the
body text from STDIN.

### --detached / -d

The `--detached` aka `-d` flag does not pull updates from the main or perennial
branch. This allows you to build out your branch stack and decide when to pull
in changes from other developers.

### --dry-run

Use the `--dry-run` flag allows to test-drive this command. It prints the Git
commands that would be run but doesn't execute them.

### --title / -t

When called with the `--title <title>` aka `-t` flag, the _propose_ command
pre-populate the title of the pull request to the given text.

### --verbose / -v

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.

### Configuration

You can configure the hosting platform type with the
[hosting-platform](../preferences/hosting-platform.md) setting.

When using SSH identities, this command uses the hostname in the
[hosting-origin-hostname](../preferences/hosting-origin-hostname.md) setting.
