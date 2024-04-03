Feature: delete a local branch

  Background:
    Given the local feature branches "current" and "other"
    And the commits
      | BRANCH  | LOCATION | MESSAGE      |
      | current | local    | local commit |
    And an uncommitted file
    And the current branch is "current" and the previous branch is "other"
    When I run "git-town kill"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                        |
      | current | git fetch --prune --tags       |
      |         | git add -A                     |
      |         | git commit -m "WIP on current" |
      |         | git checkout other             |
      | other   | git branch -D current          |
    And the current branch is now "other"
    And the branches are now
      | REPOSITORY | BRANCHES    |
      | local      | main, other |
      | origin     | main        |
    And this lineage exists now
      | BRANCH | PARENT |
      | other  | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                                       |
      | other   | git branch current {{ sha 'WIP on current' }} |
      |         | git checkout current                          |
      | current | git reset --soft HEAD~1                       |
    And the current branch is now "current"
    And the uncommitted file still exists
    And the initial commits exist
    And the initial branches and lineage exist
