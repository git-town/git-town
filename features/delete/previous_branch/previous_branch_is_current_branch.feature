Feature: deleting a branch without a useful previous branch setting

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | current | feature | main   | local     |
      | other   | feature | main   | local     |
    And the commits
      | BRANCH  | LOCATION | MESSAGE        |
      | current | local    | current commit |
      | other   | local    | other commit   |
    And the current branch is "current"
    And the current branch is "current" and the previous branch is "current"
    When I run "git-town delete"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                  |
      | current | git fetch --prune --tags |
      |         | git checkout main        |
      | main    | git branch -D current    |
    And this lineage exists now
      """
      main
        other
      """
    And the branches are now
      | REPOSITORY | BRANCHES    |
      | local      | main, other |
      | origin     | main        |
    And these commits exist now
      | BRANCH | LOCATION | MESSAGE      |
      | other  | local    | other commit |
    And no uncommitted files exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                       |
      | main   | git branch current {{ sha 'current commit' }} |
      |        | git checkout current                          |
    And the initial branches and lineage exist now
    And the initial commits exist now
