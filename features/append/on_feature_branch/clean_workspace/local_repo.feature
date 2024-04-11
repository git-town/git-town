Feature: in a local repo

  Background:
    Given my repo does not have an origin
    And the current branch is a local feature branch "existing"
    And the commits
      | BRANCH   | LOCATION | MESSAGE         |
      | existing | local    | existing commit |
    When I run "git-town append new"

  Scenario: result
    Then it runs the commands
      | BRANCH   | COMMAND                  |
      | existing | git merge --no-edit main |
      |          | git checkout -b new      |
    And the current branch is now "new"
    And these commits exist now
      | BRANCH   | LOCATION | MESSAGE         |
      | existing | local    | existing commit |
      | new      | local    | existing commit |
    And this lineage exists now
      | BRANCH   | PARENT   |
      | existing | main     |
      | new      | existing |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH   | COMMAND               |
      | new      | git checkout existing |
      | existing | git branch -D new     |
    And the current branch is now "existing"
    And the initial commits exist
    And the initial branches and lineage exist
