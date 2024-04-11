Feature: append to a perennial branch

  Background:
    Given the perennial branches "qa" and "production"
    And the commits
      | BRANCH     | LOCATION | MESSAGE           |
      | production | origin   | production commit |
    And the current branch is "production"
    When I run "git-town append new"

  Scenario: result
    Then it runs the commands
      | BRANCH     | COMMAND                        |
      | production | git fetch --prune --tags       |
      |            | git rebase origin/production   |
      |            | git checkout -b new production |
    And the current branch is now "new"
    And these commits exist now
      | BRANCH     | LOCATION      | MESSAGE           |
      | new        | local         | production commit |
      | production | local, origin | production commit |
    And this lineage exists now
      | BRANCH | PARENT     |
      | new    | production |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH     | COMMAND                                     |
      | new        | git checkout production                     |
      | production | git reset --hard {{ sha 'initial commit' }} |
      |            | git branch -D new                           |
    And the current branch is now "production"
    And the initial commits exist
    And the initial lineage exists
