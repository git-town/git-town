# Releasing new versions of Git Town

This guide is for maintainers who make releases of Git Town.

### bump the version

- search-and-replace the old version with the new version
- if bumping the major version, also update `github.com/git-town/git-town/v7/`

### create a GitHub release

On a Linux machine:

- install [hub](https://github.com/github/hub#installation)
- install [goreleaser](https://goreleaser.com/install)
- update `RELEASE_NOTES.md` on `master`
- create and push a new Git Tag for the release: `git tag -a v7.7.0 -m v7.7.0`
- `env GITHUB_TOKEN=<your Github token> VERSION=7.7.0 make release-linux`
  - or omit the Github token and enter your credentials when asked
- this opens a release in the browser
- copy the release notes into it
- delete the other invalid release that the script has created

On a Windows machine, in Git Bash:

- install [hub](https://github.com/github/hub#installation)
- install [go-msi](https://github.com/mh-cbon/go-msi#install)
- install [wix](https://wixtoolset.org/releases)
- optionally install
  [.NET 3.5](https://dotnet.microsoft.com/download/dotnet-framework)
- `env VERSION=7.7.0 make msi` to create the Windows installer
- test the created Windows installer in the `dist` directory
- `env GITHUB_TOKEN=<your Github token> VERSION=7.7.0 make release-win`
- this opens the release in the browser
- verify that it added the `.msi` file
- publish the release

### create a Homebrew release

TODO: try the new `brew bump-formula-pr` command next time.

- fork [Homebrew](https://github.com/Homebrew/homebrew-core)
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
