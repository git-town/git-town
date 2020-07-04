# Releasing new versions of Git Town

This guide is for maintainers who make releases of Git Town.

### create a GitHub release

On a Windows machine, in Git Bash:

- install [hub](https://github.com/github/hub#installation)
- create a feature branch which updates `RELEASE_NOTES.md`
- get the feature branch reviewed and merged
- create and push a new Git Tag for the release: `git tag -a v4.0 -m v4.0`
- `env GITHUB_TOKEN=<your Github token> make release`
  - or omit the Github token and enter your credentials when asked
- this opens a release in the browser - review and publish it
- it creates two releases right now, delete the one without attached files

### create a Homebrew release

- fork [Homebrew](https://github.com/Homebrew/homebrew)
- update `Library/Formula/git-town.rb`
  - get the sha256 by downloading the release (`.tar.gz`) and using
    `shasum -a 256 /path/to/file`
  - ignore the `bottle` block, the homebrew maintainers update it
- create a pull request and get it merged

### Arch Linux

Flag the package out of date on the right hand side menu of
[Git Town's AUR page](https://aur.archlinux.org/packages/git-town/).
[allonsy](https://github.com/allonsy) will update the package.

### debugging

To test the goreleaser setup:

```
goreleaser --snapshot --skip-publish --rm-dist
```
