Feature: rename a local-only branch

  Background:
    Given the current branch is a local feature branch "old"
    And the commits
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
      | old    | local         | old commit  |
    When I run "git-town rename-branch new"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | old    | git fetch --prune --tags |
      |        | git branch new old       |
      |        | git checkout new         |
      | new    | git branch -D old        |
    And the current branch is now "new"
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
      | new    | local         | old commit  |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                               |
      | new    | git branch old {{ sha 'old commit' }} |
      |        | git checkout old                      |
      | old    | git branch -D new                     |
    And the current branch is now "old"
    And the initial branches and lineage exist
