Feature: in a local repo

  Background:
    Given a local Git repo
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | existing | feature | main   | local     |
    And the current branch is "existing"
    And the commits
      | BRANCH   | LOCATION | MESSAGE         |
      | existing | local    | existing commit |
    When I run "git-town append new"

  Scenario: result
    Then it runs the commands
      | BRANCH   | COMMAND                       |
      | existing | git merge --no-edit --ff main |
      |          | git checkout -b new           |
    And the current branch is now "new"
    And these commits exist now
      | BRANCH   | LOCATION | MESSAGE         |
      | existing | local    | existing commit |
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
    And the initial commits exist now
    And the initial branches and lineage exist now
