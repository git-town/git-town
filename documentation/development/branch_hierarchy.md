# Branch Hierarchy

Git Town divides all branches into three categories:
* main branch
* perennial branches
* feature branches


## Main branch

The main development branch from which all feature branches are cut from and merged in.
This is stored in the git config under `git-town.main-branch-name`


## Perennial branches

These are branches that serve some special purpose, such as deployment.
Perennial branches cannot be killed or shipped,
and only rebase with their own tracking branch when synced.
These are stored in the git config under `git-town.perennial-branch-names` as a space seperated list.


## Nested Feature Branches

Since code reviews can take a while,
many developers work on several features in parallel.
These features often depend on each other.
To support this common use case, Git Town provides an hierarchical branching model
that is more opinionated than the very generic branching of vanilla Git.
In Git Town's world, feature branches can be "children" of other feature branches.

As an example, lets assume a repo with the following setup:

```
-o--o-- master
  \
   o--o--o-- feature1
       \
        o-- feature2
```

In this example, feature 1 (which was cut straight from the master branch) is currently under review.
While waiting for the LGTM there, the developer has started to work on the next feature.
This work (let's call it "feature 2") needs some of the changes that are introduced by feature 1.
Since feature 1 hasn't shipped yet, we can't cut feature 2 straight off master,
but must cut it off feature 1, so that feature 2 sees the changes made by feature 1.

This means the feature branch `feature1` is cut directly from `master`,
and `feature2` is cut from `feature1`, making it a child branch of `feature1`.

This "ancestry line" of branches is preserved at all times,
and impacts a lot of Git Town's commands.
For example, branches cannot be shipped before their ancestor branches.
When syncing, Git Town syncs the parent branch first,
then merges the parent branch into its children branches.
When creating a pull request for `feature2`,
Git Town only displays the changes between `feature2` and `feature1`,
not the diff against `master`.

Git Town stores the immediate parent of each feature branch in the git config under `git-town.<branch_name>.parent`.
and the full ancestral line, top-down, as a space seperated list under `git-town-branch.<branch_name>.ancestors`

For this example it would store
```
git-town-branch.feature1.parent=master
git-town-branch.feature1.ancestors=master

git-town-branch.feature2.parent=feature1
git-town-branch.feature2.ancestors=master feature1
```
