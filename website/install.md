Git Town runs anywhere Git and [Bash](https://www.gnu.org/software/bash/bash.html) run.


<table>
  <tr>
    <th width="300px">
      Using <a href="http://brew.sh">Homebrew</a>
    </th>
    <th width="400px">
      Manually
    </th>
    <th width="300px">
      Debian Based Linux
    </th>
  </tr>
  <tr class="subhead">
    <td colspan="3">
      <b>Install</b>
    </td>
  </tr>
  <tr>
    <td>
      <code>brew install git-town</code>
    </td>
    <td>
      <ul>
        <li>clone the repo to your machine (into DIR)</li>
        <li>add DIR/src to your <code>$PATH</code></li>
        <li>add DIR/man to your <code>$MANPATH</code></li>
      </ul>
    </td>
    <td>
      <ul>
        <li>Download the latest deb file from the <a href="https://github.com/Originate/git-town/releases">releases page</a></li>
        <li>Either open the deb file in your favorite gui package manager or run <code>dpkg -i $PATH_TO_DEBFILE</code></li>
      </ul>
    </td>
  </tr>
  <tr class="subhead">
    <td colspan="3">
      <b>Update</b>
    </td>
  </tr>
  <tr>
    <td>
      <code>brew update</code><br>
      <code>brew upgrade git-town</code>
    </td>
    <td>
      <ul>
        <li>run <code>git pull</code> in DIR</li>
      </ul>
    </td>
    <td>
      <ul>
        <li>redownload a fresh deb file and run the same instructions as above</li>
      </ul>
    </td>
  </tr>
  <tr class="subhead">
    <td colspan="3">
      <b>Uninstall</b>
    </td>
  </tr>
  <tr>
    <td>
      <code>brew uninstall git-town</code><br>
    </td>
    <td>
      <ul>
        <li>remove DIR</li>
        <li>remove DIR/src from your <code>$PATH</code></li>
        <li>remove DIR/man from your <code>$MANPATH</code></li>
      </ul>
    </td>
    <td>
      <ul>
        <li>run <code>apt-get remove gittown</code></li>
      </ul>
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
