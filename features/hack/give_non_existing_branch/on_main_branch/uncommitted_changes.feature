Feature: on the main branch with uncommitted changes

  Background:
    Given a Git repo with origin
    And the commits
      | BRANCH | LOCATION | MESSAGE     |
      | main   | origin   | main commit |
    And the current branch is "main"
    When I run "git-town hack new"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND             |
      | main   | git checkout -b new |
    And the current branch is now "new"
    And the initial commits exist now
    And this lineage exists now
      | BRANCH | PARENT |
      | new    | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND           |
      | new    | git checkout main |
      | main   | git branch -D new |
    And the current branch is now "main"
    And the initial commits exist now
    And the initial branches and lineage exist now
