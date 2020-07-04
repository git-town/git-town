Git Town runs anywhere Git runs.

# macOS

There are several options to install Git Town on macOS:

1. using [Homebrew](https://brew.sh):
   - Install: run `brew install git-town`
   - Update: run `brew update` and then run `brew upgrade git-town`
   - Uninstall: run `brew uninstall git-town`

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

# Debian/Ubuntu based distros

- Install / Update:
  - Download the deb file from the latest release
    [here](https://github.com/git-town/git-town/releases).
  - Run `dpkg -i /path/to/debfile`
- Uninstall: run `apt-get remove gittown`

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
