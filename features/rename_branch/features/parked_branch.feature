Feature: rename a parked branch

  Background:
    Given a Git repo with origin
    And the branch
      | NAME   | TYPE   | PARENT | LOCATIONS     |
      | parked | parked | main   | local, origin |
    And the current branch is "parked"
    And the commits
      | BRANCH | LOCATION      | MESSAGE             |
      | parked | local, origin | experimental commit |
    When I run "git-town rename-branch parked new"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | parked | git fetch --prune --tags |
      |        | git branch new parked    |
      |        | git checkout new         |
      | new    | git push -u origin new   |
      |        | git push origin :parked  |
      |        | git branch -D parked     |
    And the current branch is now "new"
    And the parked branches are now "new"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE             |
      | new    | local, origin | experimental commit |
    And this lineage exists now
      | BRANCH | PARENT |
      | new    | main   |

  Scenario: undo
    Given I ran "git-town rename-branch --force parked new"
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                           |
      | new    | git branch parked {{ sha 'experimental commit' }} |
      |        | git push -u origin parked                         |
      |        | git push origin :new                              |
      |        | git checkout parked                               |
      | parked | git branch -D new                                 |
    And the current branch is now "parked"
    And the parked branches are now "parked"
    And the initial commits exist
    And the initial branches and lineage exist
