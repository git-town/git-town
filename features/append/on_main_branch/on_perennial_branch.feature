Feature: append to a perennial branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME       | TYPE      | PARENT | LOCATIONS     |
      | qa         | perennial |        | local, origin |
      | production | perennial |        |               |
    And the commits
      | BRANCH     | LOCATION | MESSAGE           |
      | production | origin   | production commit |
    And the current branch is "production"
    When I run "git-town append new"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH     | COMMAND                                       |
      | production | git fetch --prune --tags                      |
      |            | git rebase origin/production --no-update-refs |
      |            | git checkout -b new                           |
    And the current branch is now "new"
    And these commits exist now
      | BRANCH     | LOCATION      | MESSAGE           |
      | production | local, origin | production commit |
    And this lineage exists now
      | BRANCH | PARENT     |
      | new    | production |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH     | COMMAND                                     |
      | new        | git checkout production                     |
      | production | git reset --hard {{ sha 'initial commit' }} |
      |            | git branch -D new                           |
    And the current branch is now "production"
    And the initial commits exist now
    And the initial lineage exists now
