# hosting-platform

```
git-town.hosting-platform=<github|gitlab|bitbucket|gitea>
```

To talk to the API of your code hosting platform, Git Town needs to know which
code hosting platform (GitHub, Gitlab, Bitbucket, etc) you use. Git Town can
automatically figure out the code hosting platform by looking at the URL of the
`origin` remote. In cases where that's not successful, for example when using
private instances of code hosting platforms, you can tell Git Town which code
hosting platform you use via the _hosting-platform_ preference. To set it, run

```
git config [--global] git-town.hosting-platform <name>
```

The optional `--global` flag applies this setting to all Git repositories on
your local machine. When not present, the setting applies to the current repo.
`<name>` can be "github", "gitlab", "gitea", or "bitbucket".
