Feature: deleting a branch without a useful previous branch setting

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | current | feature | main   | local     |
      | other   | feature | main   | local     |
    And the current branch is "current"
    And the commits
      | BRANCH  | LOCATION | MESSAGE        |
      | current | local    | current commit |
      | other   | local    | other commit   |
    And the current branch is "current" and the previous branch is "current"
    And an uncommitted file
    When I run "git-town delete"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                          |
      | current | git fetch --prune --tags                         |
      |         | git add -A                                       |
      |         | git commit -m "Committing WIP on deleted branch" |
      |         | git checkout main                                |
      | main    | git branch -D current                            |
    And the current branch is now "main"
    And no uncommitted files exist now
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
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                         |
      | main    | git branch current {{ sha 'Committing WIP on deleted branch' }} |
      |         | git checkout current                                            |
      | current | git reset --soft HEAD~1                                         |
    And the current branch is now "current"
    And the uncommitted file still exists
    And the initial commits exist now
    And the initial branches and lineage exist now
