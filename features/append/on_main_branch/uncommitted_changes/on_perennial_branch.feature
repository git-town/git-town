Feature: append to a perennial branch

  Background:
    Given a Git repo clone
    And the branches
      | NAME       | TYPE      | PARENT | LOCATIONS     |
      | qa         | perennial |        | local, origin |
      | production | perennial |        |               |
    And the commits
      | BRANCH     | LOCATION | MESSAGE           |
      | production | origin   | production commit |
    And the current branch is "production"
    And an uncommitted file
    When I run "git-town append new"

  Scenario: result
    Then it runs the commands
      | BRANCH     | COMMAND             |
      | production | git add -A          |
      |            | git stash           |
      |            | git checkout -b new |
      | new        | git stash pop       |
    And the current branch is now "new"
    And the initial commits exist
    And this lineage exists now
      | BRANCH | PARENT     |
      | new    | production |
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH     | COMMAND                 |
      | new        | git add -A              |
      |            | git stash               |
      |            | git checkout production |
      | production | git branch -D new       |
      |            | git stash pop           |
    And the current branch is now "production"
    And the initial commits exist
    And the initial lineage exists
    And the uncommitted file still exists
