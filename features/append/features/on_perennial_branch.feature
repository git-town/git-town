Feature: append to a perennial branch

  Background:
    Given my repo has the perennial branches "qa" and "production"
    And my repo contains the commits
      | BRANCH     | LOCATION | MESSAGE           |
      | production | remote   | production_commit |
    And I am on the "production" branch
    When I run "git-town append new-child"

  Scenario: result
    Then it runs the commands
      | BRANCH     | COMMAND                         |
      | production | git fetch --prune --tags        |
      |            | git rebase origin/production    |
      |            | git branch new-child production |
      |            | git checkout new-child          |
    And I am now on the "new-child" branch
    And my repo now has the commits
      | BRANCH     | LOCATION      | MESSAGE           |
      | new-child  | local         | production_commit |
      | production | local, remote | production_commit |
    And Git Town is now aware of this branch hierarchy
      | BRANCH    | PARENT     |
      | new-child | production |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH     | COMMAND                 |
      | new-child  | git checkout production |
      | production | git branch -D new-child |
    And I am now on the "production" branch
    And my repo now has the commits
      | BRANCH     | LOCATION      | MESSAGE           |
      | production | local, remote | production_commit |
    And Git Town now has the original branch hierarchy
