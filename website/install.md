Git Town runs anywhere Git runs.

# macOS

There are several options to install Git Town on macOS:

1. using [Homebrew](https://brew.sh):

   - Install: run `brew install git-town`
   - Update: run `brew update` and then run `brew upgrade git-town`
   - Uninstall: run `brew uninstall git-town`

2. manually:
   - download an
     [archive containing the binary](https://github.com/git-town/git-town/releases/download/v7.4.0/git-town_7.4.0_macOS_intel_64.tar.gz)
   - extract with `tar -xzf git-town_7.4.0_macOS_intel_64.tar.gz`
   - move the `git-town` binary into a directory listed in your `$PATH`, for
     example `/usr/local/bin`

---

# Windows

There are several options to install Git Town on Windows:

1. download and run the [Windows installer]()

2. install using [scoop](https://scoop.sh):

   ```
   scoop bucket add org https://github.com/git-town/scoop.git
   scoop install git-town
   ```

3. install manually: download the latest binary and put it somewhere into your
   `%PATH%`.

---

# Linux

There are several options to install Git Town on Linux distributions:

1. [Homebrew](https://brew.sh): `brew install git-town`

2. using your package manager:

   - on Debian-based systems, download the
     [.deb](https://github.com/git-town/git-town/releases/download/v7.4.0/git-town_7.4.0_linux_intel_64.deb)
     file and run `sudo apt-get install git-town_7.4.0_linux_intel_64.deb`
   - on RedHat-based systems download the
     [.rpm](https://github.com/git-town/git-town/releases/download/v7.4.0/git-town_7.4.0_linux_intel_64.rpm)
     file and run `rpm -i git-town_7.4.0_linux_intel_64.rpm`

3. manually
   - download an
     [archive containing the binary](https://github.com/git-town/git-town/releases/download/v7.4.0/git-town_7.4.0_linux_intel_64.tar.gz)
   - extract with `tar -xzf git-town_7.4.0_linux_intel_64.tar.gz`
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

# Manual

- Install: download the Git Town binary for your platform from the
  [releases page](https://github.com/git-town/git-town/releases), rename it to
  `git-town`, make it executable with `chmod +x git-town` and put it somewhere
  in your `$PATH`
- Update: install the new version over the old version
- Uninstall:
  - remove Git Town's aliases: `git town alias false`
  - remove the Git Town configuration from your repositories: in each repo, run
    `git town config reset`
  - delete the `git-town` binary from your hard drive

---

#### Install autocompletion

- follow instructions given by `git-town help completions`

#### Notifications about new releases

- Subscribe to our <a href="https://github.com/git-town/git-town/releases.atom">
  release feed <i class="ion-social-rss accent-color"></i></a> to never miss a
  new release! If you prefer email notifications, please try
  [sibbell.com](https://sibbell.com).
