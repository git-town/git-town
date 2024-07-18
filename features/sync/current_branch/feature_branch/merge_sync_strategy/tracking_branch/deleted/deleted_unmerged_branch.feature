Feature: sync a branch with unmerged commits whose tracking branch was deleted

  Background:
    Given a Git repo clone
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | branch-1 | feature | main   | local, origin |
      | branch-2 | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         |
      | branch-1 | local, origin | branch-1 commit |
      | branch-2 | local, origin | branch-2 commit |
    And origin deletes the "branch-2" branch
    And the current branch is "branch-2"
    And an uncommitted file
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH   | COMMAND                       |
      | branch-2 | git fetch --prune --tags      |
      |          | git add -A                    |
      |          | git stash                     |
      |          | git checkout main             |
      | main     | git rebase origin/main        |
      |          | git checkout branch-2         |
      | branch-2 | git merge --no-edit --ff main |
      |          | git stash pop                 |
    And it prints:
      """
      Branch "branch-2" was deleted at the remote but the local branch contains unshipped changes.
      """
    And the current branch is now "branch-2"
    And the uncommitted file still exists
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE         |
      | branch-1 | local, origin | branch-1 commit |
      | branch-2 | local         | branch-2 commit |
    And the initial branches and lineage exist

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH   | COMMAND       |
      | branch-2 | git add -A    |
      |          | git stash     |
      |          | git stash pop |
    And the current branch is now "branch-2"
    And the uncommitted file still exists
    And the initial branches and lineage exist
