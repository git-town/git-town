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
      | old    | git merge --no-edit main |
      |        | git stash pop            |
    And it prints:
      """
      Branch "old" was deleted at the remote but the local branch contains unshipped changes.
      """
    And the current branch is now "old"
    And the uncommitted file still exists
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE       |
      | active | local, origin | active commit |
      | old    | local         | old commit    |
    And the initial branches and lineage exist

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND       |
      | old    | git add -A    |
      |        | git stash     |
      |        | git stash pop |
    And the current branch is now "old"
    And the uncommitted file still exists
    And the initial branches and lineage exist
