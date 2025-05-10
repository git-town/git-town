@smoke
Feature: prepend a prototype branch to a feature branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | old  | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE    |
      | old    | local, origin | old commit |
    And the current branch is "old"
    When I run "git-town prepend --prototype parent"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                     |
      | old    | git fetch --prune --tags    |
      |        | git checkout -b parent main |
    And branch "parent" now has type "prototype"
    And branch "old" still has type "feature"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE    |
      | old    | local, origin | old commit |
    And this lineage exists now
      | BRANCH | PARENT |
      | old    | parent |
      | parent | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND              |
      | parent | git checkout old     |
      | old    | git branch -D parent |
    And branch "old" still has type "feature"
    And the initial commits exist now
    And the initial lineage exists now
