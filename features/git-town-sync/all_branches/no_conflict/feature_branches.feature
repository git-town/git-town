Feature: git-town sync --all: syncs all feature branches

  Background:
    Given my repository has the feature branches "feature-1" and "feature-2"
    And the following commits exist in my repository
      | BRANCH    | LOCATION         | MESSAGE          | FILE NAME     |
      | main      | remote           | main commit      | main_file     |
      | feature-1 | local and remote | feature-1 commit | feature1_file |
      | feature-2 | local and remote | feature-2 commit | feature2_file |
    And I am on the "feature-1" branch
    And my workspace has an uncommitted file
    When I run `git-town sync --all`


  Scenario: result
    Then Git Town runs the commands
      | BRANCH    | COMMAND                              |
      | feature-1 | git fetch --prune                    |
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
    And my repository has the following commits
      | BRANCH    | LOCATION         | MESSAGE                            | FILE NAME     |
      | main      | local and remote | main commit                        | main_file     |
      | feature-1 | local and remote | feature-1 commit                   | feature1_file |
      |           |                  | main commit                        | main_file     |
      |           |                  | Merge branch 'main' into feature-1 |               |
      | feature-2 | local and remote | feature-2 commit                   | feature2_file |
      |           |                  | main commit                        | main_file     |
      |           |                  | Merge branch 'main' into feature-2 |               |
