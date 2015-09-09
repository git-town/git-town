Git Town runs anywhere Git and [Bash](https://www.gnu.org/software/bash/bash.html) run.


<table>
  <tr>
    <th width="300px">
      Using <a href="http://brew.sh">Homebrew</a>
    </th>
    <th width="400px">
      Manually
    </th>
  </tr>
  <tr class="subhead">
    <td colspan="2">
      <b>Install</b>
    </td>
  </tr>
  <tr>
    <td>
      <code>brew tap Originate/gittown</code><br>
      <code>brew install git-town</code>
    </td>
    <td>
      <ul>
        <li>clone the repo to your machine (into DIR)</li>
        <li>add DIR/src to your <code>$PATH</code></li>
        <li>add DIR/man to your <code>$MANPATH</code></li>
        <li>
          install <a href="http://en.wikipedia.org/wiki/Dialog_(software)">Dialog</a>
          (used by <a href="/documentation/commands/git-extract.md">git extract</a>)
        </li>
      </ul>
    </td>
  </tr>
  <tr class="subhead">
    <td colspan="2">
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
  </tr>
  <tr class="subhead">
    <td colspan="2">
      <b>Uninstall</b>
    </td>
  </tr>
  <tr>
    <td>
      <code>brew uninstall git-town</code><br>
      <code>brew untap Originate/gittown</code>
    </td>
    <td>
      <ul>
        <li>remove DIR</li>
        <li>remove DIR/src from your <code>$PATH</code></li>
        <li>remove DIR/man from your <code>$MANPATH</code></li>
      </ul>
    </td>
  </tr>
</table>


#### Install autocompletion

* for [Fish shell](http://fishshell.com): `git town install-fish-autocompletion`


#### Notifications for new releases

* Subscribe to our
  [release RSS feed](https://github.com/Originate/git-town/releases.atom)
  to never miss a new release!
  If you prefer email notifications, please try https://sibbell.com.
