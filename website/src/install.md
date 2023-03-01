# Installation

Git Town ships as a single self-contained binary. It doesn't bundle a Git client
but uses the native Git installation on your machine.

[![Packaging status](https://repology.org/badge/vertical-allrepos/git-town.svg)](https://repology.org/project/git-town/versions)

## macOS

The easiest way to install Git Town on macOS is via Homebrew:

```
brew install git-town
```

Git Town is available via [MacPorts](https://ports.macports.org/port/git-town).
You can also [install Git Town manually](#manual-installation) or
[compile from source](#compile-from-source).

## Windows

The easiest way to install Git Town on Windows is to
[download the Windows installer](https://github.com/git-town/git-town/releases/download/v7.9.0/git-town_7.9.0_windows_intel_64.msi)
and run it. You can also install Git Town via [scoop](https://scoop.sh):

```
scoop bucket add org https://github.com/git-town/scoop.git
scoop install git-town
```

You can also [install Git Town manually](#manual-installation) or
[compile from source](#compile-from-source).

If you use the Windows Subsystem for Linux, please install
[wsl-open](https://www.npmjs.com/package/wsl-open) to allow the commands
[git town repo](https://git-town.com/commands/repo.md) and
[git town new-pull-request](https://git-town.com/commands/new-pull-request.md)
to open a browser window for you.

## Linux

On Debian-based systems,
[download](https://github.com/git-town/git-town/releases/latest) the `.deb` file
matching your CPU architecture and run:

```
sudo apt-get install git-town_7.9.0_linux_intel_64.deb
```

On RedHat-based systems
[download](https://github.com/git-town/git-town/releases/latest) the `.rpm` file
matching your CPU architecture and run

```
rpm -i git-town_7.9.0_linux_intel_64.rpm
```

On Arch Linux, install the
[git-town](https://aur.archlinux.org/packages/git-town) package from the AUR.

You can install Git Town via
[Homebrew for Linux](https://docs.brew.sh/Homebrew-on-Linux):

```
brew install git-town
```

You can also [install Git Town manually](#manual-installation) or
[compile from source](#compile-from-source).

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

Subscribe to our
[release feed](https://github.com/git-town/git-town/releases.atom) to never miss
a new release!

## Uninstall

To remove Git Town from your system:

1. remove Git Town's aliases: `git town aliases false`
2. remove the Git Town configuration from your repositories: in each repo, run
   `git town config reset`
3. uninstall the program or manually delete the binary
