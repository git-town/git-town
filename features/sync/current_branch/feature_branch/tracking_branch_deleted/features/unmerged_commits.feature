Feature: sync a branch with unmerged commits whose tracking branch was deleted

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

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | old    | git fetch --prune --tags |
      |        | git add -A               |
      |        | git stash                |
      |        | git checkout main        |
      | main   | git rebase origin/main   |
      |        | git checkout old         |
      | old    | git stash pop            |
    And it prints:
      """
      The branch "old" was deleted on the remote but the local branch on this machine contains unshipped changes.
      I am therefore not removing this branch. Run "git town diff-parent" to see the changes.
      """
    And the current branch is now "old"
    And the uncommitted file still exists
    And the branches are now
      | REPOSITORY | BRANCHES          |
      | local      | main, active, old |
      | origin     | main, active      |

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
