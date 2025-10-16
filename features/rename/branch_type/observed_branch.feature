Feature: rename an observed branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE     | PARENT | LOCATIONS     |
      | observed | observed | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE               |
      | observed | local, origin | somebody elses commit |
    And the current branch is "observed"
    When I run "git-town rename observed new"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                        |
      | observed | git fetch --prune --tags       |
      |          | git branch --move observed new |
      |          | git checkout new               |
      | new      | git push -u origin new         |
      |          | git push origin :observed      |
    And this lineage exists now
      """
      main
        new
      """
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE               |
      | new    | local, origin | somebody elses commit |
    And branch "new" still has type "observed"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                               |
      | new      | git branch observed {{ sha 'somebody elses commit' }} |
      |          | git push -u origin observed                           |
      |          | git checkout observed                                 |
      | observed | git branch -D new                                     |
      |          | git push origin :new                                  |
    And the initial branches and lineage exist now
    And branch "observed" still has type "observed"
    And the initial commits exist now
