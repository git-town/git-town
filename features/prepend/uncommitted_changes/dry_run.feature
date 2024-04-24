Feature: dry-run prepending a branch to a feature branch

  Background:
    Given the current branch is a feature branch "old"
    And the commits
      | BRANCH | LOCATION      | MESSAGE    |
      | old    | local, origin | old commit |
    And an uncommitted file
    When I run "git-town prepend parent --dry-run"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                        |
      | old    | git add -A                     |
      |        | git stash                      |
      |        | git checkout main              |
      | main   | git rebase origin/main         |
      |        | git checkout old               |
      | old    | git merge --no-edit origin/old |
      |        | git merge --no-edit --ff main  |
      |        | git checkout -b parent main    |
      | parent | git stash pop                  |
    And the current branch is still "old"
    And the uncommitted file still exists
    And the initial commits exist
    And the initial branches and lineage exist

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "old"
    And the uncommitted file still exists
    And the initial commits exist
    And the initial lineage exists
