Feature: delete the current parked branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | parked  | parked  | main   | local, origin |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
      | parked  | local, origin | parked commit  |
    And the current branch is "parked"
    And the current branch is "parked" and the previous branch is "feature"
    When I run "git-town delete"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                  |
      | parked  | git fetch --prune --tags |
      |         | git push origin :parked  |
      |         | git checkout feature     |
      | feature | git branch -D parked     |
    And this lineage exists now
      """
      main
        feature
      """
    And the branches are now
      | REPOSITORY    | BRANCHES      |
      | local, origin | main, feature |
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
    And no uncommitted files exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                     |
      | feature | git branch parked {{ sha 'parked commit' }} |
      |         | git push -u origin parked                   |
      |         | git checkout parked                         |
    And the initial branches and lineage exist now
    And branch "parked" now has type "parked"
    And the initial commits exist now
