# Releasing new versions of Git Town

This guide is for maintainers who make releases of Git Town.

- in a branch:
  - update [CHANGELOG.md](CHANGELOG.md)
  - search-and-replace the old version with the new version
    - triple-digits: `10.0.2`
    - double-digits: `10.0`
    - don't change existing version numbers in CHANGELOG.md
  - if bumping the major version:
    - update `github.com/git-town/git-town/v10/` everywhere in this repo
    - update `github.com/git-town/git-town/v10` (without trailing slash)
- ship the branch
- `git tag v10.0.2 && git push --tags``
- open the new release in the browser
- copy the changelog into the release notes
- publish the release
- merge the `main` branch into the `public` branch

### update the website

The website publishes from the `public` branch to avoid listing unreleased
features to the public.

```
git checkout public
git merge main
git push
```

### debugging the release script

Debugging is best done on a separate copy of this codebase. This avoids
accidental releases, which then trigger other automation at the Homebrew servers
etc.

The release script is in PowerShell because creating the .msi file requires
Windows. So it's best to debug it on a Windows machine.

Test the compile step:

```powershell
.\run-that-app goreleaser@1.22.1 --snapshot --skip-publish --clean
```

Test the entire release pipeline:

- this MUST happen in a separate repo

```powershell
git push ; git tag -d v0.0.1 ; git push origin :v0.0.1 ; git tag v0.0.1 ; git push --tags
$env:GITHUB_TOKEN="<github token>"; $env:VERSION="0.0.1"; $env:TODAY="today"; .\tools\release.ps1
```

### release platforms

- HomeBrew: Git Town is in the auto-updating
  [core formulae](https://formulae.brew.sh/formula/git-town)
- Scoop: Git Town is in the auto-updating
  [core manifests](https://github.com/ScoopInstaller/Main/blob/master/bucket/git-town.json)
- Arch Linux: the [AUR package](https://aur.archlinux.org/packages/git-town)
  auto-updates
