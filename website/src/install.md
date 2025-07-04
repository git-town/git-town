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
  [Git Town Windows installer](https://github.com/git-town/git-town/releases/download/v18.0.0/git-town_windows_intel_64.msi)

If you use the Windows Subsystem for Linux, please install
[wsl-open](https://www.npmjs.com/package/wsl-open) to allow the commands
[git town repo](https://www.git-town.com/commands/repo) and
[git town propose](https://www.git-town.com/commands/propose) to open a browser
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
[git-town](https://archlinux.org/packages/extra/x86_64/git-town/) package from
official package repositories:

```
pacman -S git-town
```

On openSUSE Tumbleweed, install the
[git-town](https://build.opensuse.org/package/show/openSUSE:Factory/git-town)
package from the official package repositories:

```
sudo zypper in git-town
```

OpenSUSE 15.x and 16.x users can use these steps to install Git Town:

```bash
# replace 15.6 with 16.0 in the following command
zypper ar -f -r https://download.opensuse.org/repositories/home:/ojkastl_buildservice:/git-town/15.6/home:ojkastl_buildservice:git-town.repo
zypper refresh # accept the GPG key for the devel:tools:scm repository
zypper install git-town
```

There are separate packages for the shell completions called
`git-town-bash-completion`, `git-town-zsh-completion`, and
`git-town-fish-completion`.

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

## Manual installation

```
curl https://www.git-town.com/install.sh | sh
```

For a fully custom installation,
[download](https://github.com/git-town/git-town/releases) the archive matching
your CPU architecture, extract it, and move the `git-town` executable into a
directory listed in your `$PATH`, for example `/usr/local/bin`.

## Compile from source

If you have the [Go compiler](https://go.dev) installed, you can compile the
latest version of Git Town from source by running:

```
go get github.com/git-town/git-town/v21
```

## New releases

Subscribe to the
[Git Town release feed](https://github.com/git-town/git-town/releases.atom) to
get notifications about new releases.

## Uninstall

To remove Git Town from your system:

1. Remove the Git Town configuration from your repositories: in each repo, run
   `git town config remove`
2. If your operating system or package manager provides an uninstaller for Git
   Town, run it. If you installed Git Town manually, delete the binary.
