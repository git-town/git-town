# Install

Git Town runs anywhere Git runs.

### macOS

There are several options to install Git Town on macOS. The easiest way is via Homebrew:

```
brew install git-town
```

You can also install Git Town manually:

- download an archive containing the binary
- extract with `tar -xzf git-town_7.6.0_macOS_intel_64.tar.gz`
- move the git-town binary into a directory listed in your `$PATH`, for example `/usr/local/bin`

### Windows

There are several options to install Git Town on Windows. The easiest way is to [download the Windows installer](https://github.com/git-town/git-town/releases/download/v7.6.0/git-town_7.6.0_windows_intel_64.msi) and run it.
You can also install Git Town via scoop:

```
scoop bucket add org https://github.com/git-town/scoop.git
scoop install git-town
```

As usual, you can manually [download the Windows binary](https://github.com/git-town/git-town/releases/download/v7.6.0/git-town_7.6.0_windows_intel_64.zip) and put it somewhere into your `%PATH%`.

If you use the Windows Subsystem for Linux, install [wsl-open](https://www.npmjs.com/package/wsl-open) to make the commands [git town repo](https://github.com/git-town/git-town/blob/master/documentation/commands/repo.md) and [git town new-pull-request](https://github.com/git-town/git-town/blob/master/documentation/commands/new-pull-request.md) work.

### Linux

There are several options to install Git Town on Linux distributions. On Debian-based systems, [download](https://github.com/git-town/git-town/releases/latest) the `.deb` file matching your CPU architecture and run:

```
sudo apt-get install git-town_7.6.0_linux_intel_64.deb
```

On RedHat-based systems [download](https://github.com/git-town/git-town/releases/latest) the `.rpm` file matching your CPU architecture and run

```
rpm -i git-town_7.6.0_linux_intel_64.rpm
```

You can install Git Town via [Homebrew for Linux](https://docs.brew.sh/Homebrew-on-Linux):

```
brew install git-town
```

Or install it manually:

- [download](https://github.com/git-town/git-town/releases/latest) the archive containing the Linux binary
- extract: `tar -xzf git-town_7.6.0_linux_intel_64.tar.gz`
- move the git-town binary into a directory listed in your `$PATH`, for example `/usr/local/bin`

Arch Linux

    Install / Update: install git-town from the AUR
        You may use your favorite AUR helper
        Alternatively, you can install manually from the AUR by:
            downloading the latest tarball from this page
            unpacking the tarball
            cd into the newly unpacked directory, and running makepkg -si
    Uninstall: run pacman -R git-town

Install autocompletion

    follow instructions given by git-town help completions

Notifications about new releases

    Subscribe to our release feed to never miss a new release!

Uninstall

When uninstalling Git Town:

    remove Git Town's aliases: git town alias false
    remove the Git Town configuration from your repositories: in each repo, run git town config reset
    delete the git-town binary from your hard drive
