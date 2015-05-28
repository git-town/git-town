This _code hosting_ interface makes functionality on
third-party source code hosting services like [GitHub](http://github.com)
or [Bitbucket](http://bitbucket.org) available.

The particular driver to be used is automatically determined by analyzing
the `origin` remote in the Git repo that the command is run.


## Interface

<table>
  <thead>
    <tr></tr>
    <tr>
      <th width="345px">function</th>
      <th>description</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>
        <code>create_pull_request &lt;repository&gt; &lt;branch&gt;</code>
      </td>
      <td>
        opens a browser window that allows to create a pull request
        for the given branch in the given repository against the main branch
      </td>
    </tr>
    <tr>
      <td>
        <code>show_repo &lt;repository&gt;</code>
      </td>
      <td>
        opens a browser window showing the home page of the given repository
        on its hosting platform
      </td>
    </tr>
  </tbody>
</table>
