Feature: offline mode

  Background:
    Given a Git repo with origin
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | old  | feature | main   | local, origin |
    And offline mode is enabled
    And the current branch is "old"
    And the commits
      | BRANCH | LOCATION      | MESSAGE    |
      | old    | local, origin | old commit |
    When I run "git-town prepend new"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                 |
      | old    | git checkout main                       |
      | main   | git rebase origin/main --no-update-refs |
      |        | git checkout old                        |
      | old    | git merge --no-edit --ff main           |
      |        | git merge --no-edit --ff origin/old     |
      |        | git checkout -b new main                |
    And the current branch is now "new"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE    |
      | old    | local, origin | old commit |
    And this lineage exists now
      | BRANCH | PARENT |
      | new    | main   |
      | old    | new    |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND           |
      | new    | git checkout old  |
      | old    | git branch -D new |
    And the current branch is now "old"
    And the initial commits exist now
    And the initial lineage exists now