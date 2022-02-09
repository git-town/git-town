Feature: prepend a branch to a feature branch

  Background:
    Given the current branch is a feature branch "old"
    And the commits
      | BRANCH | LOCATION      | MESSAGE    |
      | old    | local, origin | old commit |
    And an uncommitted file
    When I run "git-town prepend parent"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | old    | git fetch --prune --tags |
      |        | git add -A               |
      |        | git stash                |
      |        | git checkout main        |
      | main   | git rebase origin/main   |
      |        | git branch parent main   |
      |        | git checkout parent      |
      | parent | git stash pop            |
    And the current branch is now "parent"
    And the uncommitted file still exists
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE    |
      | old    | local, origin | old commit |
    And this branch hierarchy exists now
      | BRANCH | PARENT |
      | old    | parent |
      | parent | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND              |
      | parent | git add -A           |
      |        | git stash            |
      |        | git checkout main    |
      | main   | git branch -d parent |
      |        | git checkout old     |
      | old    | git stash pop        |
    And the current branch is now "old"
    And the uncommitted file still exists
    And now the initial commits exist
    And the initial branch hierarchy exists
