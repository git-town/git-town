Git Town runs anywhere Git runs.

# macOS

There are several options to install Git Town on macOS:

1. using [Homebrew](https://brew.sh): `brew install git-town`

2. manually:
   - download an
     [archive containing the binary](https://github.com/git-town/git-town/releases/download/v7.5.0/git-town_7.5.0_macOS_intel_64.tar.gz)
   - extract with `tar -xzf git-town_7.5.0_macOS_intel_64.tar.gz`
   - move the `git-town` binary into a directory listed in your `$PATH`, for
     example `/usr/local/bin`

---

# Windows

There are several options to install Git Town on Windows:

1. Windows installer:
   [download](https://github.com/git-town/git-town/releases/download/v7.5.0/git-town_7.5.0_windows_intel_64.msi)

2. [scoop](https://scoop.sh):

   ```
   scoop bucket add org https://github.com/git-town/scoop.git
   scoop install git-town
   ```

3. manually: download the
   [ZIP file](https://github.com/git-town/git-town/releases/download/v7.5.0/git-town_7.5.0_windows_intel_64.zip),
   and put the binary inside it somewhere into your `%PATH%`.

If you use the Windows Subsystem for Linux, install
[wsl-open](https://www.npmjs.com/package/wsl-open) to make the commands
[git town repo](https://github.com/git-town/git-town/blob/master/documentation/commands/repo.md)
and
[git town new-pull-request](https://github.com/git-town/git-town/blob/master/documentation/commands/new-pull-request.md)
work.

---

# Linux

There are several options to install Git Town on Linux distributions:

1. [Homebrew](https://brew.sh): `brew install git-town`

2. using your package manager:

   - on Debian-based systems, download the
     [.deb](https://github.com/git-town/git-town/releases/download/v7.5.0/git-town_7.5.0_linux_intel_64.deb)
     file and run `sudo apt-get install git-town_7.5.0_linux_intel_64.deb`
   - on RedHat-based systems download the
     [.rpm](https://github.com/git-town/git-town/releases/download/v7.5.0/git-town_7.5.0_linux_intel_64.rpm)
     file and run `rpm -i git-town_7.5.0_linux_intel_64.rpm`

3. manually
   - download an
     [archive](https://github.com/git-town/git-town/releases/download/v7.5.0/git-town_7.5.0_linux_intel_64.tar.gz)
     containing the binary
   - extract with `tar -xzf git-town_7.5.0_linux_intel_64.tar.gz`
   - move the `git-town` binary into a directory listed in your `$PATH`, for
     example `/usr/local/bin`

---

# Arch Linux

- Install / Update: install
  [git-town](https://aur.archlinux.org/packages/git-town/) from the AUR
  - You may use your favorite AUR helper
  - Alternatively, you can install manually from the AUR by:
    - downloading the latest tarball
      [from this page](https://aur.archlinux.org/packages/git-town/)
    - unpacking the tarball
    - `cd` into the newly unpacked directory, and running `makepkg -si`
- Uninstall: run `pacman -R git-town`

---

# Install autocompletion

- follow instructions given by `git-town help completions`

# Notifications about new releases

- Subscribe to our <a href="https://github.com/git-town/git-town/releases.atom">
  release feed <i class="ion-social-rss accent-color"></i></a> to never miss a
  new release!

# Uninstall

When uninstalling Git Town:

- remove Git Town's aliases: `git town alias false`
- remove the Git Town configuration from your repositories: in each repo, run
  `git town config reset`
- delete the `git-town` binary from your hard drive
