Feature: delete the current contribution branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME         | TYPE         | PARENT | LOCATIONS     |
      | contribution | contribution |        | local, origin |
      | feature      | feature      | main   | local, origin |
    And the commits
      | BRANCH       | LOCATION      | MESSAGE             |
      | contribution | local, origin | contribution commit |
      | feature      | local, origin | feature commit      |
    And the current branch is "contribution"
    And the current branch is "contribution" and the previous branch is "feature"
    When I run "git-town delete"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH       | COMMAND                    |
      | contribution | git fetch --prune --tags   |
      |              | git checkout feature       |
      | feature      | git branch -D contribution |
    And this lineage exists now
      """
      main
        feature
      """
    And the branches are now
      | REPOSITORY | BRANCHES                    |
      | local      | main, feature               |
      | origin     | main, contribution, feature |
    And these commits exist now
      | BRANCH       | LOCATION      | MESSAGE             |
      | feature      | local, origin | feature commit      |
      | contribution | origin        | contribution commit |
    And no uncommitted files exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                 |
      | feature | git branch contribution {{ sha 'contribution commit' }} |
      |         | git checkout contribution                               |
    And the initial branches and lineage exist now
    And branch "contribution" now has type "contribution"
    And the initial commits exist now
