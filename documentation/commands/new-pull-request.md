#### NAME

new-pull-request - create a new pull request


#### SYNOPSIS

```
git town new-pull-request
```


#### DESCRIPTION

Syncs the current branch
and opens a browser window to the new pull request page of your repository.

The form is pre-populated for the current branch
so that the pull request only shows the changes made
against the immediate parent branch.

Supported only for repositories hosted on [GitHub](http://github.com/),
[GitLab](http://gitlab.com/), and [Bitbucket](https://bitbucket.org/).
When using hosted versions of GitHub, GitLab, or Bitbucket,
make sure that your SSH identity contains the phrase "github", "gitlab" or
"bitbucket", so that Git Town can derive which hosting service you use.
Example: your SSH identity should be something like `git@github-as-account1:Originate/git town.git`
