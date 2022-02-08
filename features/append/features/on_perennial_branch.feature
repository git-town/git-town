Feature: append to a perennial branch

  Background:
    Given my repo has the perennial branches "qa" and "production"
    And the commits
      | BRANCH     | LOCATION | MESSAGE           |
      | production | origin   | production commit |
    And I am on the "production" branch
    When I run "git-town append new"

  Scenario: result
    Then it runs the commands
      | BRANCH     | COMMAND                      |
      | production | git fetch --prune --tags     |
      |            | git rebase origin/production |
      |            | git branch new production    |
      |            | git checkout new             |
    And I am now on the "new" branch
    And now these commits exist
      | BRANCH     | LOCATION      | MESSAGE           |
      | new        | local         | production commit |
      | production | local, origin | production commit |
    And Git Town is now aware of this branch hierarchy
      | BRANCH | PARENT     |
      | new    | production |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH     | COMMAND                 |
      | new        | git checkout production |
      | production | git branch -D new       |
    And I am now on the "production" branch
    And now these commits exist
      | BRANCH     | LOCATION      | MESSAGE           |
      | production | local, origin | production commit |
    And Git Town is now aware of the initial branch hierarchy
