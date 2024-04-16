Feature: killing a branch without a useful previous branch setting

  Background:
    Given the current branch is a local feature branch "current"
    And a local feature branch "other"
    And the commits
      | BRANCH  | LOCATION | MESSAGE        |
      | current | local    | current commit |
      | other   | local    | other commit   |
    And the current branch is "current" and the previous branch is "current"
    And an uncommitted file
    When I run "git-town kill"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                        |
      | current | git fetch --prune --tags       |
      |         | git add -A                     |
      |         | git commit -m "WIP on current" |
      |         | git checkout main              |
      | main    | git branch -D current          |
    And the current branch is now "main"
    And no uncommitted files exist
    And the branches are now
      | REPOSITORY | BRANCHES    |
      | local      | main, other |
      | origin     | main        |
    And these commits exist now
      | BRANCH | LOCATION | MESSAGE      |
      | other  | local    | other commit |
    And this lineage exists now
      | BRANCH | PARENT |
      | other  | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                                       |
      | main    | git branch current {{ sha 'WIP on current' }} |
      |         | git checkout current                          |
      | current | git reset --soft HEAD~1                       |
    And the current branch is now "current"
    And the uncommitted file still exists
    And the initial commits exist
    And the initial branches and lineage exist
