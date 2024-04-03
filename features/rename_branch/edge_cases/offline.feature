Feature: offline mode

  Background:
    Given offline mode is enabled
    And the current branch is a feature branch "old"
    And the commits
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
      | old    | local, origin | old commit  |
    When I run "git-town rename-branch new"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND            |
      | old    | git branch new old |
      |        | git checkout new   |
      | new    | git branch -D old  |
    And the current branch is now "new"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
      | new    | local         | old commit  |
      | old    | origin        | old commit  |
    And this lineage exists now
      | BRANCH | PARENT |
      | new    | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                               |
      | new    | git branch old {{ sha 'old commit' }} |
      |        | git checkout old                      |
      | old    | git branch -D new                     |
    And the current branch is now "old"
    And the initial commits exist
    And the initial branches and lineage exist
