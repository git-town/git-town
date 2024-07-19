@smoke
Feature: rename the current branch

  Background:
    Given a Git repo clone
    And the branch
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | old  | feature | main   | local, origin |
    And the current branch is "old"
    And the commits
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
      | old    | local, origin | old commit  |
    When I run "git-town rename-branch new"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | old    | git fetch --prune --tags |
      |        | git branch new old       |
      |        | git checkout new         |
      | new    | git push -u origin new   |
      |        | git push origin :old     |
      |        | git branch -D old        |
    And the current branch is now "new"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
      | new    | local, origin | old commit  |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                               |
      | new    | git branch old {{ sha 'old commit' }} |
      |        | git push -u origin old                |
      |        | git push origin :new                  |
      |        | git checkout old                      |
      | old    | git branch -D new                     |
    And the current branch is now "old"
    And the initial branches and lineage exist
