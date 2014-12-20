#### NAME

git-town - general Git Town help, view and change Git Town configuration


#### SYNOPSIS

```
git town
git town config
git town help
git town main-branch [<branchname>]
git town non-feature-branches [(--add | --remove) <branchname>]
git town version
```

#### COMMANDS

* *help*
> View help screen. Running `git town` bare will also show the help screen.

* *version*
> View the Git Town version.

* *config*
> View your current Git Town configuration.

* *main-branch*
> View your main-branch configuration.
>
> With an optional branch name `<branchname>`, specify a branch to assign as the main branch.

* *non-feature-branches*
> View your non-feature branch configuration.
>
> With the `--add` or `--remove` option paired with a `<branchname>`, you may update your non-feature branches accordingly.