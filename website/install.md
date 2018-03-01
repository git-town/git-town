Git Town runs anywhere Git runs.

# Mac (using [Homebrew](http://brew.sh))

* Install: run `brew install git-town`
* Update: run `brew update` and then run `brew upgrade git-town`
* Uninstall: run `brew uninstall git-town`

_Note: `4.0.0` dropped support for Yosemite because we now specify a minimum Git version of 2.7.0,
which is higher than the default version installed on that OS.
If you would like to use Git Town on Yosemite,
please use the manual install and ensure your Git version is 2.7.0 or higher._

---

# Arch Linux

* Install / Update: install [git-town](https://aur.archlinux.org/packages/git-town/) from the AUR
  * You may use your favorite AUR helper
  * Alternatively, you can install manually from the AUR by:
    * downloading the latest tarball [from this page](https://aur.archlinux.org/packages/git-town/)
    * unpacking the tarball
    * `cd` into the newly unpacked directory, and running `makepkg -si`
* Uninstall: run `pacman -R git-town`

---

# Debian/Ubuntu based distros

* Install / Update:
  * Download the deb file from the latest release [here](https://github.com/Originate/git-town/releases).
  * Run `dpkg -i /path/to/debfile`
* Uninstall: run `apt-get remove git-town`

---

# Manual

* Install: download the Git Town binary for your platform from the
  [releases page](https://github.com/Originate/git-town/releases),
  rename it to `git-town`, make it executable with `chmod +x git-town`
  and put it somewhere in your `$PATH`
* Update: install the new version over the old version
* Uninstall:
  * remove Git Town's aliases: `git town alias false`
  * remove the Git Town configuration from your repositories:
    in each repo, run `git town config reset`
  * delete the `git-town` binary from your hard drive

---

#### Install autocompletion

* for [Fish shell](http://fishshell.com): `git town install-fish-autocompletion`

#### Notifications about new releases

* Subscribe to our
  <a href="https://github.com/Originate/git-town/releases.atom">
  release feed <i class="ion-social-rss accent-color"></i></a> to never miss a new release!
  If you prefer email notifications, please try [sibbell.com](https://sibbell.com).
