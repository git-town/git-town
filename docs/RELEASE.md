# Releasing new versions of Git Town

This guide is for maintainers who make releases of Git Town.

- in a branch:
  - update CHANGELOG.md
  - run `make stats-release` and copy the release stats and contributors into
    CHANGELOG.md
  - verify that all newly added Git Town commands are not hidden and listed in
    the README file
  - verify that the website content reflects all the changes made
  - search-and-replace the old version with the new version
    - triple-digits: `16.7.0`
    - double-digits: `16.7`
    - its best to do this manually so that you don't change unrelated version
      numbers in CHANGELOG.md, lockfiles, and data files
  - if bumping the major version:
    - update `github.com/git-town/git-town/v16/` everywhere in this repo
    - update `github.com/git-town/git-town/v16` (without trailing slash)
- ship the branch
- `git sync --all && git checkout public && git merge main && git push && git checkout main && git branch -d public`
- wait for Netlify to publish the website:
  https://app.netlify.com/sites/git-town/deploys
- in a new branch:
  - update the changelog with links to the website
  - ship the branch
- `git sync --all && git checkout main && git tag v16.7.0 && git push --tags`
- wait for CI to finish the release:
  https://github.com/git-town/git-town/actions
- edit the new unpublished release:
  https://github.com/git-town/git-town/releases
- copy the changelog into the release notes
- publish the release

### update the website

The website publishes from the `public` branch to avoid listing unreleased
features to the public.

```
git checkout public
git merge main
git push
```

### debugging the release script

Debugging is best done on a separate fork of this codebase. This avoids
accidental releases to the official location, which then trigger other
automation at Homebrew, Scoop, Arch Linux, etc.

The release script is written in PowerShell because creating the .msi file
requires Windows. It's best to debug it on a Windows machine.

Test the compile step:

```powershell
.\rta goreleaser@1.22.1 --snapshot --skip-publish --clean
```

Test the entire release pipeline:

- this MUST happen in a separate repo

```powershell
git push ; git tag -d v0.0.1 ; git push origin :v0.0.1 ; git tag v0.0.1 ; git push --tags
$env:GITHUB_TOKEN="<github token>"; $env:VERSION="0.0.1"; $env:TODAY="today"; .\tools\release.ps1
```

### performing a manual release

If the release script fails in production and doesn't create the release, and/or
you want to investigate some release code, you can perform the release manually
on a Windows machine using PowerShell.

- move the affected Git tag to HEAD but only locally, DON'T PUSH THE TAG TO
  ORIGIN

  ```
  git tag -d <tag>
  git tag <tag>
  ```

- run the release script
  ```powershell
  $env:GITHUB_TOKEN="<token>"
  $env:CHOCOLATEY_API_KEY="<key>"
  .\tools\release.ps1
  ```

- delete the local Git tag and download the real one from origin

  ```
  git tag -d <tag>
  git fetch --tags
  ```

### release platforms

- HomeBrew: Git Town is in the auto-updating
  [core formulae](https://formulae.brew.sh/formula/git-town)
- Scoop: Git Town is in the auto-updating
  [core manifests](https://github.com/ScoopInstaller/Main/blob/master/bucket/git-town.json)
- Arch Linux: the [AUR package](https://aur.archlinux.org/packages/git-town)
  auto-updates
