Feature: rename a parked branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE   | PARENT | LOCATIONS     |
      | parked | parked | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE             |
      | parked | local, origin | low-priority commit |
    And the current branch is "parked"
    When I run "git-town rename parked new"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                      |
      | parked | git fetch --prune --tags     |
      |        | git branch --move parked new |
      |        | git checkout new             |
      | new    | git push -u origin new       |
      |        | git push origin :parked      |
    And this lineage exists now
      """
      main
        new
      """
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE             |
      | new    | local, origin | low-priority commit |
    And branch "new" still has type "parked"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                           |
      | new    | git branch parked {{ sha 'low-priority commit' }} |
      |        | git push -u origin parked                         |
      |        | git checkout parked                               |
      | parked | git branch -D new                                 |
      |        | git push origin :new                              |
    And the initial branches and lineage exist now
    And branch "parked" still has type "parked"
