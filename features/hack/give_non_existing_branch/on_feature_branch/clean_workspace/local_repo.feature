Feature: local repo

  Background:
    Given the current branch is a feature branch "existing"
    And my repo does not have an origin
    And the commits
      | BRANCH | LOCATION | MESSAGE     |
      | main   | local    | main commit |
    When I run "git-town hack new"

  Scenario: result
    Then it runs the commands
      | BRANCH   | COMMAND                  |
      | existing | git checkout -b new main |
    And the current branch is now "new"
    And these commits exist now
      | BRANCH | LOCATION | MESSAGE     |
      | main   | local    | main commit |
      | new    | local    | main commit |
    And this lineage exists now
      | BRANCH   | PARENT |
      | existing | main   |
      | new      | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH   | COMMAND               |
      | new      | git checkout existing |
      | existing | git branch -D new     |
    And the current branch is now "existing"
    And the initial commits exist
    And this lineage exists now
      | BRANCH   | PARENT |
      | existing | main   |
