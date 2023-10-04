Feature: sync a branch with unmerged commits whose tracking branch was deleted

  # TODO: decide what to do here
  #
  # Option A: The branch was deleted on the remote, so it should be deleted locally as well.
  # This is especially true in this example where the local client doesn't contain any additional changes
  # beyond those that existed on origin and were deleted there.
  #
  # Option B: If the branch truly contains local-only changes that were not on origin
  # when origin deleted the branch, then it should not be deleted locally.
  # It's just hard to determine that since Git doesn't give us a SHA for the now deleted remote branch.
  # We might be able to look this up in the Git history, but there doesn't seem to be a straightforward way
  # since the history doesn't show old branches.

  Background:
    Given the feature branches "active" and "old"
    And the commits
      | BRANCH | LOCATION      | MESSAGE       |
      | active | local, origin | active commit |
      | old    | local, origin | old commit    |
    And origin deletes the "old" branch
    And the current branch is "old"
    And an uncommitted file
    When I run "git-town sync"

  @this
  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | old    | git fetch --prune --tags |
      |        | git add -A               |
      |        | git stash                |
      |        | git checkout main        |
      | main   | git rebase origin/main   |
      |        | git branch -D old        |
      |        | git stash pop            |
    And the current branch is now "main"
    And the uncommitted file still exists
    And the branches are now
      | REPOSITORY    | BRANCHES     |
      | local, origin | main, active |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND           |
      | old    | git add -A        |
      |        | git stash         |
      |        | git checkout main |
      | main   | git checkout old  |
      | old    | git stash pop     |
    And the current branch is now "old"
    And the uncommitted file still exists
    And the initial branches and hierarchy exist
