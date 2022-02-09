Feature: in a local repo

  Background:
    Given my repo does not have an origin
    And a local feature branch "existing"
    And the commits
      | BRANCH   | LOCATION | MESSAGE         |
      | existing | local    | existing commit |
    And the current branch is "existing"
    And my workspace has an uncommitted file
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
    And my workspace still contains my uncommitted file
    And now these commits exist
      | BRANCH   | LOCATION | MESSAGE         |
      | existing | local    | existing commit |
    And Git Town is now aware of this branch hierarchy
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
      | existing | git branch -d new     |
      |          | git stash pop         |
    And the current branch is now "existing"
    And now the initial commits exist
    And my workspace still contains my uncommitted file
    And the initial branches and hierarchy exist
