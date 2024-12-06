Feature: delete a parent branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE    | PARENT | LOCATIONS     |
      | feature-1 | feature | main   | local, origin |
    And the commits
      | BRANCH    | LOCATION      | MESSAGE          |
      | feature-1 | local, origin | feature-1 commit |
    And the branches
      | NAME      | TYPE    | PARENT    | LOCATIONS     |
      | feature-2 | feature | feature-1 | local, origin |
    And the commits
      | BRANCH    | LOCATION      | MESSAGE          |
      | feature-2 | local, origin | feature-2 commit |
    And the branches
      | NAME      | TYPE    | PARENT    | LOCATIONS     |
      | feature-3 | feature | feature-2 | local, origin |
    And the commits
      | BRANCH    | LOCATION      | MESSAGE          |
      | feature-3 | local, origin | feature-3 commit |
    And the current branch is "feature-3"
    And an uncommitted file
    When I run "git-town delete feature-2"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH    | COMMAND                    |
      | feature-3 | git fetch --prune --tags   |
      |           | git add -A                 |
      |           | git stash                  |
      |           | git push origin :feature-2 |
      |           | git branch -D feature-2    |
      |           | git stash pop              |
    And the current branch is now "feature-3"
    And the uncommitted file still exists
    And the branches are now
      | REPOSITORY    | BRANCHES                   |
      | local, origin | main, feature-1, feature-3 |
    And these commits exist now
      | BRANCH    | LOCATION      | MESSAGE          |
      | feature-1 | local, origin | feature-1 commit |
      | feature-3 | local, origin | feature-2 commit |
      |           |               | feature-3 commit |
    And this lineage exists now
      | BRANCH    | PARENT    |
      | feature-1 | main      |
      | feature-3 | feature-1 |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                           |
      | feature-3 | git add -A                                        |
      |           | git stash                                         |
      |           | git branch feature-2 {{ sha 'feature-2 commit' }} |
      |           | git push -u origin feature-2                      |
      |           | git stash pop                                     |
    And the current branch is now "feature-3"
    And the uncommitted file still exists
    And the initial commits exist now
    And the initial branches and lineage exist now
