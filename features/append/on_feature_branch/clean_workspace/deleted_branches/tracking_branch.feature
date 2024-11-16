Feature: append a branch to a branch whose tracking branch was deleted

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | shipped | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | shipped | local, origin | shipped commit |
    And origin ships the "shipped" branch using the squash-merge ship-strategy
    And the current branch is "shipped"
    When I run "git-town append new"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                 |
      | shipped | git fetch --prune --tags                |
      |         | git checkout main                       |
      | main    | git rebase origin/main --no-update-refs |
      |         | git branch -D shipped                   |
      |         | git checkout -b new                     |
    And Git Town prints:
      """
      deleted branch "shipped"
      """
    And the current branch is now "new"
    And the branches are now
      | REPOSITORY | BRANCHES  |
      | local      | main, new |
      | origin     | main      |
    And this lineage exists now
      | BRANCH | PARENT |
      | new    | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                       |
      | new     | git checkout main                             |
      | main    | git reset --hard {{ sha 'initial commit' }}   |
      |         | git branch shipped {{ sha 'shipped commit' }} |
      |         | git checkout shipped                          |
      | shipped | git branch -D new                             |
    And the current branch is now "shipped"
    And the initial branches and lineage exist now
