# code-hosting-platform

```
git-town.code-hosting-platform=<github|gitlab|bitbucket|gitea>
```

To talk to the API of your code hosting service, Git Town needs to know which
code hosting service (GitHub, Gitlab, Bitbucket, etc) you use. Git Town can
automatically figure out the code hosting platform by looking at the URL of the
`origin` remote. In cases where that's not successful, for example when using
private instances of code hosting services, you can tell Git Town which code
hosting service you use via the _code-hosting-platform_ preference. To set it,
run

```
git config [--global] git-town.code-hosting-platform <name>
```

The optional `--global` flag applies this setting to all Git repositories on
your local machine. When not present, the setting applies to the current repo.
`<name>` can be "github", "gitlab", "gitea", or "bitbucket".
