Feature: prepend a branch to a feature branch in a dirty workspace using the "compress" sync strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | old  | feature | main   | local, origin |
    And the current branch is "old"
    And the commits
      | BRANCH | LOCATION      | MESSAGE    |
      | old    | local, origin | old commit |
    And Git Town setting "sync-feature-strategy" is "compress"
    And an uncommitted file
    When I run "git-town prepend parent"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                     |
      | old    | git add -A                  |
      |        | git stash                   |
      |        | git checkout -b parent main |
      | parent | git stash pop               |
    And the current branch is now "parent"
    And the uncommitted file still exists
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE    |
      | old    | local, origin | old commit |
    And this lineage exists now
      | BRANCH | PARENT |
      | old    | parent |
      | parent | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND              |
      | parent | git add -A           |
      |        | git stash            |
      |        | git checkout old     |
      | old    | git branch -D parent |
      |        | git stash pop        |
    And the current branch is now "old"
    And the uncommitted file still exists
    And the initial commits exist
    And the initial lineage exists
