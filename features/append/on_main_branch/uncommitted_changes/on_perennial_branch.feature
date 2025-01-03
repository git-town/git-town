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
      | BRANCH     | COMMAND                     |
      | production | git add -A                  |
      |            | git stash -m "Git Town WIP" |
      |            | git checkout -b new         |
      | new        | git stash pop               |
    And the current branch is now "new"
    And the initial commits exist now
    And this lineage exists now
      | BRANCH | PARENT     |
      | new    | production |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH     | COMMAND                     |
      | new        | git add -A                  |
      |            | git stash -m "Git Town WIP" |
      |            | git checkout production     |
      | production | git branch -D new           |
      |            | git stash pop               |
    And the current branch is now "production"
    And the initial commits exist now
    And the initial lineage exists now
