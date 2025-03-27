Feature: hoisting a branch out of a stack

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | branch-1 | feature | main   | local     |
    And the commits
      | BRANCH   | LOCATION | MESSAGE  |
      | branch-1 | local    | commit 1 |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS |
      | branch-2 | feature | branch-1 | local     |
    And the commits
      | BRANCH   | LOCATION | MESSAGE  |
      | branch-2 | local    | commit 2 |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS |
      | branch-3 | feature | branch-2 | local     |
    And the commits
      | BRANCH   | LOCATION | MESSAGE  |
      | branch-3 | local    | commit 3 |
    And the current branch is "branch-2"
    When I run "git rebase --onto main branch-1 branch-2"
    And I run "git checkout branch-3"
    When I run "git rebase --onto branch-1 branch-2 branch-3"
  # When I run "git-town hoist"

  @this
  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND |
    # | branch-2 | git fetch --prune --tags                |
    # |          | git checkout main                       |
    # | main     | git rebase origin/main --no-update-refs |
    # |          | git checkout old                        |
    # | old      | git merge --no-edit --ff main           |
    # |          | git merge --no-edit --ff origin/old     |
    # |          | git checkout -b parent main             |
    # And the current branch is still "branch-2"
    And these commits exist now
      | BRANCH   | LOCATION | MESSAGE  |
      | branch-1 | local    | commit 1 |
      | branch-2 | local    | commit 2 |
      | branch-3 | local    | commit 1 |
      |          |          | commit 3 |

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
