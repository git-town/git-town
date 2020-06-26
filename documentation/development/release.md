# Releasing new versions of Git Town

This guide is for maintainers only.

### create a GitHub release

- create a feature branch which updates `RELEASE_NOTES.md`
- get the feature branch reviewed and merged
- create and push a new Git Tag for the release
  - `git tag -m release -a v4.0`
  - `git push --tags`
- run `goreleaser`
- review the release and publish it

### create a Homebrew release

- fork [Homebrew](https://github.com/Homebrew/homebrew)
- update `Library/Formula/git-town.rb`
  - get the sha256 by downloading the release (`.tar.gz`) and using
    `shasum -a 256 /path/to/file`
  - ignore the `bottle` block. It is updated by the homebrew maintainers
- create a pull request and get it merged

### debugging

- test the goreleaser setup: `goreleaser --snapshot --skip-publish --rm-dist`
