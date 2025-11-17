# git town propose

<a type="command-summary">

```command-summary
git town propose [-b <text> | --body <text>] [-f <path> | --body-file <path>] [-t <text> | --title <text>] [-s | --stack] [--auto-resolve] [--dry-run] [-v | --verbose] [-h | --help]
```

</a>

The _propose_ command helps create a new pull request (also known as merge
request) for the current feature branch. It opens your forge's website to create
a new proposal in your browser and pre-populates information like branch and
source/target repository. It also [syncs](sync.md) the branch to merge before
opening the pull request in [detached](sync.md#-d--detached--no-detached) mode.

Proposing prototype and parked branches makes them feature branches.

You can create pull requests for repositories hosted on:

- [Bitbucket](https://bitbucket.org)
- [Forgejo](https://forgejo.org)
- [Gitea](https://gitea.com)
- [GitHub](https://github.com)
- [GitLab](https://gitlab.com)

You can configure the browser which Git Town opens using the
[BROWSER environment variable](../preferences/browser.md).

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

#### `-s`<br>`--stack`

The `--stack` aka `-s` parameter makes Git Town propose all branches in the
stack that the current branch belongs to.

#### `--dry-run`

Use the `--dry-run` flag to test-drive this command. It prints the Git commands
that would be run but doesn't execute them.

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.

#### `--auto-resolve`

Disables automatic resolution of
[phantom merge conflicts](../stacked-changes.md#avoid-phantom-conflicts).

## Configuration

You can configure the forge type with the
[hosting-platform](../preferences/forge-type.md) setting.

When using SSH identities, this command uses the hostname in the
[hosting-origin-hostname](../preferences/hosting-origin-hostname.md) setting.

## See also

- [repo](repo.md) opens the website for the repository in the browser
- [ship](ship.md) ships the current branch
