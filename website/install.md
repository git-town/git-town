Git Town runs anywhere Git and [Bash](https://www.gnu.org/software/bash/bash.html) run.

# Mac (using [Homebrew](http://brew.sh))
* Install: run `brew install git-town`
* Update: run `brew update` and then run `brew upgrade git-town`
* Uninstall: run `brew uninstall git-town`

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
* Install
  * clone the repo to your machine (into DIR)
  * add DIR/src to your `$PATH`
  * add DIR/man to your `$MANPATH`
* Update: run `git pull` in DIR
* Uninstall:
  * remove DIR
  * remove DIR/src from your `$PATH`
  * remove DIR/man from your `$MANPATH`

---


#### Install autocompletion

* for [Fish shell](http://fishshell.com): `git town install-fish-autocompletion`


#### Notifications about new releases

* Subscribe to our
  <a href="https://github.com/Originate/git-town/releases.atom">
  release feed <i class="ion-social-rss accent-color"></i></a> to never miss a new release!
  If you prefer email notifications, please try [sibbell.com](https://sibbell.com).
