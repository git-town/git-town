Git Town runs anywhere Git and [Bash](https://www.gnu.org/software/bash/bash.html) run.


<table>
  <tr>
    <th width="400px">
      Install
    </th>
    <th width="400px">
      Update
    </th>
    <th width="400px">
      Uninstall
    </th>
  </tr>

  <tr class="subhead">
    <td colspan="3">
      <b>Using <a href="http://brew.sh">Homebrew</b>
    </td>
  </tr>
  <tr>
    <td>
      <li>
        Run <code>brew install git-town</code>
      </li>
    </td>
    <td>
      <li>
        Run <code>brew update</code><br> and then run
        <code>brew upgrade git-town</code>
      </li>
    </td>
    <td>
      <li>
        Run <code>brew uninstall git-town</code>
      </li>
    </td>
  </tr>

  <tr class="subhead">
    <td colspan="3">
      <b>On Arch Linux</b>
    </td>
  </tr>
  <tr>
    <td>
      <li>Install <code><a href="https://aur.archlinux.org/packages/git-town/">git-town</a></code> from the AUR
      </li>
      <li> Your may use your favorite AUR helper
      </li>
      <li>
        Alternatively, you can install manually from the AUR by:
        <ul>
          <li>downloading the tarball <a href="https://aur.archlinux.org/packages/git-town/">from this page</a></li>
          <li>unpacking the tarball </li>
          <li><code>cd</code> into the newly unpacked directory, and running <code>makepkg -si</code> </li>
        </ul>
      </li>
    </td>
    <td>
      <li> Download the newest tarball from the <a href="https://aur.archlinux.org/packages/git-town/">AUR</a> (or use your helper to reinstall <code>git-town</code>) </li>
      <li>Rerun the install instructions</li>
    </td>
    <td>
      <li>Run <code>pacman -R git-town </code></li>
    </td>
  </tr>

  <tr class="subhead">
    <td colspan="3">
      <b>On Debian based distros</b>
    </td>
  </tr>
  <tr>
    <td>
        <li>
          Download the deb file from the latest release <a href="https://github.com/Originate/git-town/releases">here</a>.
        </li>
        <li>Run <code>dpkg -i /path/to/debfile</code></li>
    </td>
    <td>
      <li>Redownload the newest deb file from the <a href="https://github.com/Originate/git-town/releases">releases page</a></li>
      <li>Rerun the install instructions</li>
    </td>
    <td>
      <li>Run <code>apt-get remove git-town</code></li>
    </td>
  </tr>

  <tr class="subhead">
    <td colspan="3">
      <b>Manually</b>
    </td>
  </tr>
  <tr>
    <td>
      <li>clone the repo to your machine (into DIR)</li>
      <li>add DIR/src to your <code>$PATH</code></li>
      <li>add DIR/man to your <code>$MANPATH</code></li>
    </td>
    <td>
      <li>run <code>git pull</code> in DIR</li>
    </td>
    <td>
      <li>remove DIR</li>
      <li>remove DIR/src from your <code>$PATH</code></li>
      <li>remove DIR/man from your <code>$MANPATH</code></li>
    </td>
  </tr>
</table>


#### Install autocompletion

* for [Fish shell](http://fishshell.com): `git town install-fish-autocompletion`


#### Notifications about new releases

* Subscribe to our
  <a href="https://github.com/Originate/git-town/releases.atom">
  release feed <i class="ion-social-rss accent-color"></i></a> to never miss a new release!
  If you prefer email notifications, please try [sibbell.com](https://sibbell.com).
