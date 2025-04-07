# git town propose

```command-summary
git town propose [-b <text> | --body <text>] [-f <path> | --body-file <path>] [-t <text> | --title <text>] [-d | --detached] [--dry-run] [-v | --verbose]
```

The _propose_ command helps create a new pull request (also known as merge
request) for the current feature branch. It opens your forge's website to create
a new proposal in your browser and pre-populates information like branch and
source/target repository. It also [syncs](sync.md) the branch to merge before
opening the pull request.

You can create pull requests for repositories hosted on:

- [Bitbucket](https://bitbucket.org)
- [Codeberg](https://codeberg.org)
- [Gitea](https://gitea.com)
- [GitHub](https://github.com)
- [GitLab](https://gitlab.com)

On non-Windows systems, Git Town will first read the `BROWSER` environment
variable to determine the browser command. If it isn't set, Git Town will try
various common commands like `open`, `xdg-open`, or `x-www-browser`.

## Options

#### `-b <text>`<br>`--body <text>`

Pre-populate the body of the pull request with the given text.

#### `-f <path>`<br>`--body-file <path>`

When called with the `--body-file` aka `-f` flag, it pre-populates the body of
the pull request with the content of the given file. The filename `-` reads the
body text from STDIN.

#### `-t <text>`<br>`--title <text>`

When called with the `--title <title>` aka `-t` flag, the _propose_ command
pre-populate the title of the pull request to the given text.

#### `-d`<br>`--detached`

The `--detached` aka `-d` flag does not pull updates from the main or perennial
branch. This allows you to build out your stack and decide when to pull in
changes from other developers.

#### `--dry-run`

Use the `--dry-run` flag to test-drive this command. It prints the Git commands
that would be run but doesn't execute them.

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.

## Configuration

You can configure the forge type with the
[hosting-platform](../preferences/forge-type.md) setting.

When using SSH identities, this command uses the hostname in the
[hosting-origin-hostname](../preferences/hosting-origin-hostname.md) setting.
