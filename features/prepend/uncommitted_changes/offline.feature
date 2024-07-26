Feature: offline mode

  Background:
    Given a Git repo with origin
    And the branch
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | old  | feature | main   | local, origin |
    And offline mode is enabled
    And the current branch is "old"
    And the commits
      | BRANCH | LOCATION      | MESSAGE    |
      | old    | local, origin | old commit |
    And an uncommitted file
    When I run "git-town prepend new"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                             |
      | old    | git add -A                          |
      |        | git stash                           |
      |        | git checkout main                   |
      | main   | git rebase origin/main              |
      |        | git checkout old                    |
      | old    | git merge --no-edit --ff origin/old |
      |        | git merge --no-edit --ff main       |
      |        | git checkout -b new main            |
      | new    | git stash pop                       |
    And the current branch is now "new"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE    |
      | old    | local, origin | old commit |
    And this lineage exists now
      | BRANCH | PARENT |
      | new    | main   |
      | old    | new    |
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND           |
      | new    | git add -A        |
      |        | git stash         |
      |        | git checkout old  |
      | old    | git branch -D new |
      |        | git stash pop     |
    And the current branch is now "old"
    And the initial commits exist
    And the initial lineage exists
    And the uncommitted file still exists
