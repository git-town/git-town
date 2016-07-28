#### NAME

git-town - general Git Town help, view and change Git Town configuration


#### SYNOPSIS

```
git town
git town config [--reset | --setup]
git town hack-push-flag [(true | false)]
git town help
git town install-fish-autocompletion
git town main-branch [<branch_name>]
git town parent-branch <child_branch_name> <parent_branch_name>
git town perennial-branches [(--add | --remove) <branch_name>]
git town pull-branch-strategy [(rebase | merge)]
git town version
```


#### OPTIONS

* *config*
> Displays the current Git Town configuration.
>
> With the `--reset` flag, cleanly remove all Git Town configuration from the current repository.
> With the `--setup` flag, start the Git Town configuration wizard.

* *hack-push-flag*
> Displays the git-hack push flag
>
> Specify a value for the git-hack push flag.
> ```bash
> git town hack-push-flag true  # (Default). Your newly-hacked branch will be pushed upon creation.
> git town hack-push-flag false # Your newly-hacked branch will not be pushed upon creation.
> ```

* *help*
> Displays the help screen.

* *install-fish-autocompletion*
> Installs the autocompletion definition for [Fish shell](http://fishshell.com)

* *main-branch*
> Displays the name of the main development branch.
>
> With an optional branch name `<branch_name>`, specify a branch to assign as the main branch.
> ```bash
> # Set "master" as the main branch
> git town main-branch master
> ````

* *perennial-branches*
> Displays the names of all perennial branches.
>
> With the `--add` or `--remove` option, you may update your perennial branches.
> ```bash
> # Register the "qa" branch as a perennial branch
> git town perennial-branches --add qa
>
> # Remove "qa" branch from the list of perennial branches
> git town perennial-branches --remove qa
> ```

* *pull-branch-strategy*
> Displays the pull branch strategy
>
> Specify a strategy to set the pull branch strategy.
> ```bash
> git town pull-branch-strategy rebase # (Default).
> git town pull-branch-strategy merge
> ```

* *set-parent-branch*
> Update the parent branch of a feature branch
>
> ```bash
> # Set the parent branch of "feature-a" to "feature-b"
> git town parent-branch feature-a feature-b
> ``

* *version*
> Displays the Git Town version.
