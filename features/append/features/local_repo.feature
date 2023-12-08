Feature: in a local repo

  Background:
    Given my repo does not have an origin
    And the current branch is a local feature branch "existing"
    And the commits
      | BRANCH   | LOCATION | MESSAGE         |
      | existing | local    | existing commit |
    And an uncommitted file
    When I run "git-town hack new"

  Scenario: result
    Then it runs the commands
      | BRANCH   | COMMAND             |
      | existing | git add -A          |
      |          | git stash           |
      |          | git branch new main |
      |          | git checkout new    |
      | new      | git stash pop       |
    And the current branch is now "new"
    And the uncommitted file still exists
    And these commits exist now
      | BRANCH   | LOCATION | MESSAGE         |
      | existing | local    | existing commit |
    And this branch lineage exists now
      | BRANCH   | PARENT |
      | existing | main   |
      | new      | main   |

  Scenario: undo
    When I run "git town undo"
    Then it runs the commands
      | BRANCH   | COMMAND               |
      | new      | git add -A            |
      |          | git stash             |
      |          | git checkout existing |
      | existing | git branch -D new     |
      |          | git stash pop         |
    And the current branch is now "existing"
    And the initial commits exist
    And the uncommitted file still exists
    And the initial branches and lineage exist
