Feature: in a local repo

  Background:
    Given my repo has a feature branch "existing"
    And my repo does not have a remote origin
    And my repo contains the commits
      | BRANCH   | LOCATION | MESSAGE         |
      | existing | local    | existing commit |
    And I am on the "existing" branch
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
    And I am now on the "new" branch
    And my workspace still contains my uncommitted file
    And my repo now has the commits
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
    And I am now on the "existing" branch
    And my repo is left with my initial commits
    And my workspace still contains my uncommitted file
    And my repo now has its initial branches and branch hierarchy
