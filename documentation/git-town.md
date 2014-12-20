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


#### OPTIONS

* *help*
> Displays the help screen.

* *version*
> Displays the Git Town version.

* *config*
> Displays the current Git Town configuration.

* *main-branch*
> Displays the name of the main development branch.
>
> With an optional branch name `<branchname>`, specify a branch to assign as the main branch.
> ```bash
> # Set "master" as the main branch
> git town main-branch master
> ```

* *non-feature-branches*
> Displays the names of all non-feature branches.
>
> With the `--add` or `--remove` option, you may update your non-feature branches.
> ```bash
> # Register the "qa" branch as a non-feature branch
> git town non-feature-branches --add qa
>
> # Remove "qa" branch from the list of non-feature branches
> git town non-feature-branches --remove qa
> ```
