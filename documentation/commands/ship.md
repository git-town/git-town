#### NAME

ship - deliver a completed feature branch

#### SYNOPSIS

```
git town ship [<branch_name>] [<commit_options>]
git town ship (--abort | --continue)
```

#### DESCRIPTION

Squash-merges the current branch, or `<branch_name>` if given,
into the main branch, resulting in linear history on the main branch.

* syncs the main branch
* pulls remote updates for `<branch_name>`
* merges the main branch into `<branch_name>`
* squash-merges `<branch_name>` into the main branch with commit message specified by the user
* pushes the main branch to the remote repository
* deletes `<branch_name>` from the local and remote repositories

Only shipping of direct children of the main branch is allowed.
To ship a nested child branch, all ancestor branches have to be shipped or killed.

##### GitHub Pull Request Integration

If you are using GitHub, this command can squash merge pull requests via the GitHub API. Setup:

1. Get a GitHub personal access token with the `repo` scope ([see how](https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line))
2. Run `git config git-town.github-token XXX` (optionally add the `--global` flag)

Now anytime you ship a branch with a pull request on GitHub, it will squash merge via the GitHub API.
It will also update the base branch for any pull requests against that branch.

#### OPTIONS

```
<branch_name>
    The branch to ship.
    If not provided, uses the current branch.

<commit_options>
    Options to pass to 'git commit' when commiting the squash-merge.

--abort
    Cancel the operation and reset the workspace to a consistent state.

--continue
    Continue the operation after resolving conflicts.
```
