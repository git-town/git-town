Feature: delete the current observed branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE     | PARENT | LOCATIONS     |
      | observed | observed |        | local, origin |
      | feature  | feature  | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         |
      | feature  | local, origin | feature commit  |
      | observed | local, origin | observed commit |
    And the current branch is "observed"
    And the current branch is "observed" and the previous branch is "feature"
    When I run "git-town delete"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                  |
      | observed | git fetch --prune --tags |
      |          | git checkout feature     |
      | feature  | git branch -D observed   |
    And this lineage exists now
      """
      main
        feature
      """
    And the branches are now
      | REPOSITORY | BRANCHES                |
      | local      | main, feature           |
      | origin     | main, feature, observed |
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE         |
      | feature  | local, origin | feature commit  |
      | observed | origin        | observed commit |
    And no uncommitted files exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                         |
      | feature | git branch observed {{ sha 'observed commit' }} |
      |         | git checkout observed                           |
    And the initial branches and lineage exist now
    And branch "observed" now has type "observed"
    And the initial commits exist now
