Feature: hoisting a branch out of a stack

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-1 | feature | main     | local, origin |
      | branch-2 | feature | branch-1 | local, origin |
      | branch-3 | feature | branch-2 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE  |
      | branch-1 | local, origin | commit 1 |
      | branch-2 | local, origin | commit 2 |
      | branch-3 | local, origin | commit 3 |
    And the current branch is "branch-2"
    When I run "git-town hoist"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                 |
      | branch-2 | git fetch --prune --tags                |
      |          | git checkout main                       |
      | main     | git rebase origin/main --no-update-refs |
      |          | git checkout old                        |
      | old      | git merge --no-edit --ff main           |
      |          | git merge --no-edit --ff origin/old     |
      |          | git checkout -b parent main             |
    And the current branch is now "parent"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE    |
      | old    | local, origin | old commit |
    And this lineage exists now
      | BRANCH | PARENT |
      | old    | parent |
      | parent | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND              |
      | parent | git checkout old     |
      | old    | git branch -D parent |
    And the current branch is now "old"
    And the initial commits exist now
    And the initial lineage exists now
