# Installation

Git Town ships as a single self-contained binary. It calls the Git executable
that is already installed on your machine.

[![Packaging status](https://repology.org/badge/vertical-allrepos/git-town.svg)](https://repology.org/project/git-town/versions)

## macOS

You can install Git Town on macOS via
[Homebrew](https://formulae.brew.sh/formula/git-town):

```
brew install git-town
```

Installation via [MacPorts](https://ports.macports.org/port/git-town) is also
supported:

```
sudo port install git-town
```

## Windows

You can install Git Town on Windows using:

- [Chocolatey](https://community.chocolatey.org/packages/git-town):
  `choco install git-town`
- [Scoop](https://scoop.sh/#/apps?q=git+town): `scoop install git-town`
- the
  [Git Town Windows installer](https://github.com/git-town/git-town/releases/download/v15.1.0/git-town_windows_intel_64.msi)

If you use the Windows Subsystem for Linux, please install
[wsl-open](https://www.npmjs.com/package/wsl-open) to allow the commands
[git town repo](https://git-town.com/commands/repo.md) and
[git town propose](https://git-town.com/commands/propose.md) to open a browser
window for you.

## Linux

On Debian-based systems,
[download](https://github.com/git-town/git-town/releases/latest) the `.deb` file
matching your CPU architecture and run:

```
sudo apt-get install git-town_linux_intel_64.deb
```

On RedHat-based systems
[download](https://github.com/git-town/git-town/releases/latest) the `.rpm` file
matching your CPU architecture and run

```
rpm -i git-town_linux_intel_64.rpm
```

On Arch Linux, install the
[git-town](https://aur.archlinux.org/packages/git-town) package from the AUR. Or
download the matching `.pkg.tar.zst` file for your architecture and run:

```
sudo pacman -U <filename>
```

You can also install Git Town on Linux via
[Homebrew for Linux](https://docs.brew.sh/Homebrew-on-Linux):

```
brew install git-town
```

You can also [install Git Town manually](#manual-installation) or
[compile from source](#compile-from-source).

## BSD

You can install Git Town on BSD via
[freshports](https://www.freshports.org/devel/git-town) or by downloading the
matching binaries from the GitHub release.

## manual installation

```
curl https://git-town.com/install.sh | sh
```

For a fully custom installation,
[download](https://github.com/git-town/git-town/releases) the archive matching
your CPU architecture, extract it, and move the `git-town` executable into a
directory listed in your `$PATH`, for example `/usr/local/bin`.

## compile from source

If you have the [Go compiler](https://go.dev) installed, you can compile the
latest version of Git Town from source by running:

```
go get github.com/git-town/git-town
```

## New releases

Subscribe to the
[Git Town release feed](https://github.com/git-town/git-town/releases.atom) to
get notifications about new releases.

## Uninstall

To remove Git Town from your system:

2. Remove the Git Town configuration from your repositories: in each repo, run
   `git town config remove`
3. If your operating system or package manager provides an uninstaller for Git
   Town, run it. If you installed Git Town manually, delete the binary.
