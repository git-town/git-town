Feature: rename a parked branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE   | PARENT | LOCATIONS     |
      | parked | parked | main   | local, origin |
    And the current branch is "parked"
    And the commits
      | BRANCH | LOCATION      | MESSAGE             |
      | parked | local, origin | low-priority commit |
    When I run "git-town rename-branch parked new"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                      |
      | parked | git fetch --prune --tags     |
      |        | git branch --move parked new |
      |        | git checkout new             |
      | new    | git push -u origin new       |
      |        | git push origin :parked      |
    And the current branch is now "new"
    And the parked branches are now "new"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE             |
      | new    | local, origin | low-priority commit |
    And this lineage exists now
      | BRANCH | PARENT |
      | new    | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                           |
      | new    | git branch parked {{ sha 'low-priority commit' }} |
      |        | git push -u origin parked                         |
      |        | git checkout parked                               |
      | parked | git branch -D new                                 |
      |        | git push origin :new                              |
    And the current branch is now "parked"
    And the parked branches are now "parked"
    And the initial commits exist now
    And the initial branches and lineage exist now
