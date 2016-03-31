#### NAME

git-town - general Git Town help, view and change Git Town configuration


#### SYNOPSIS

```
git town
git town config [--reset | --setup]
git town hack-push-strategy [(push | local)]
git town help
git town install-fish-autocompletion
git town main-branch [<branch_name>]
git town perennial-branches [(--add | --remove) <branch_name>]
git town pull-branch-strategy [(rebase | merge)]
git town version
```


#### OPTIONS

* *help*
> Displays the help screen.

* *version*
> Displays the Git Town version.

* *config*
> Displays the current Git Town configuration.
>
> With the `--reset` flag, cleanly remove all Git Town configuration from the current repository.
> With the `--setup` flag, start the Git Town configuration wizard.

* *hack-push-strategy*
> Displays the git-hack push strategy
>
> Specify a strategy to set the git-hack push strategy.
> ```bash
> git town hack-push-strategy push  # (Default). Your newly-hacked branch will be pushed upon creation.
> git town hack-push-strategy local # Your newly-hacked branch will not be pushed upon creation.
> ```

* *main-branch*
> Displays the name of the main development branch.
>
> With an optional branch name `<branch_name>`, specify a branch to assign as the main branch.
> ```bash
> # Set "master" as the main branch
> git town main-branch master
> ```

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

* *install-fish-autocompletion*
> Installs the autocompletion definition for [Fish shell](http://fishshell.com)
