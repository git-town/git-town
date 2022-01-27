Feature: git-town sync --all: syncs all feature branches

  Background:
    Given my repo has the feature branches "feature-1" and "feature-2"
    And the following commits exist in my repo
      | BRANCH    | LOCATION      | MESSAGE          |
      | main      | remote        | main commit      |
      | feature-1 | local, remote | feature-1 commit |
      | feature-2 | local, remote | feature-2 commit |
    And I am on the "feature-1" branch
    And my workspace has an uncommitted file
    When I run "git-town sync --all"

  Scenario: result
    Then it runs the commands
      | BRANCH    | COMMAND                              |
      | feature-1 | git fetch --prune --tags             |
      |           | git add -A                           |
      |           | git stash                            |
      |           | git checkout main                    |
      | main      | git rebase origin/main               |
      |           | git checkout feature-1               |
      | feature-1 | git merge --no-edit origin/feature-1 |
      |           | git merge --no-edit main             |
      |           | git push                             |
      |           | git checkout feature-2               |
      | feature-2 | git merge --no-edit origin/feature-2 |
      |           | git merge --no-edit main             |
      |           | git push                             |
      |           | git checkout feature-1               |
      | feature-1 | git push --tags                      |
      |           | git stash pop                        |
    And I am still on the "feature-1" branch
    And my workspace still contains my uncommitted file
    And all branches are now synchronized
    And my repo now has the following commits
      | BRANCH    | LOCATION      | MESSAGE                            |
      | main      | local, remote | main commit                        |
      | feature-1 | local, remote | feature-1 commit                   |
      |           |               | main commit                        |
      |           |               | Merge branch 'main' into feature-1 |
      | feature-2 | local, remote | feature-2 commit                   |
      |           |               | main commit                        |
      |           |               | Merge branch 'main' into feature-2 |
