# Releasing new versions of Git Town

This guide is for maintainers who make releases of Git Town.

### create a GitHub release

- create a feature branch which updates `RELEASE_NOTES.md`
- get the feature branch reviewed and merged
- create and push a new Git Tag for the release
  - `git tag -m release -a v4.0`
  - `git push --tags`
- run `env GITHUB_TOKEN=XXX goreleaser`
- review the release on [GitHub](https://github.com/git-town/git-town/releases)
  and publish it

### create a Homebrew release

- fork [Homebrew](https://github.com/Homebrew/homebrew)
- update `Library/Formula/git-town.rb`
  - get the sha256 by downloading the release (`.tar.gz`) and using
    `shasum -a 256 /path/to/file`
  - ignore the `bottle` block, the homebrew maintainers update it
- create a pull request and get it merged

### create the Windows installer

On a Windows machine, in Git Bash:

<pre textrun="verify-make-command">
make msi
</pre>

Manually add the generated `.msi` file to the GitHub release.

### Arch Linux

Flag the package out of date on the right hand side menu of
[Git Town's AUR page](https://aur.archlinux.org/packages/git-town/).
[allonsy](https://github.com/allonsy) will update the package.

### debugging

To test the goreleaser setup:

```
goreleaser --snapshot --skip-publish --rm-dist
```
