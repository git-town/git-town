@smoke
Feature: rename the current branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | old  | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
      | old    | local, origin | old commit  |
    And the current branch is "old"
    When I run "git-town rename new"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                   |
      | old    | git fetch --prune --tags  |
      |        | git branch --move old new |
      |        | git checkout new          |
      | new    | git push -u origin new    |
      |        | git push origin :old      |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
      | new    | local, origin | old commit  |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                               |
      | new    | git branch old {{ sha 'old commit' }} |
      |        | git push -u origin old                |
      |        | git checkout old                      |
      | old    | git branch -D new                     |
      |        | git push origin :new                  |
    And the initial branches and lineage exist now
