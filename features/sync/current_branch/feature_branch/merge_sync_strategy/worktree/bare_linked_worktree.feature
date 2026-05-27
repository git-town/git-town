Feature: sync from a linked worktree of a bare repo

  The parent branch (main) is the HEAD of a bare clone, so it appears with a
  worktreepath in git for-each-ref output. git-town must recognize that path as
  a bare repo and treat main as available — not as locked in another worktree.

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION | MESSAGE               |
      | main    | origin   | origin main commit    |
      | feature | local    | local feature commit  |
      |         | origin   | origin feature commit |
    And branch "feature" is active in a linked worktree of a bare clone
    When I run "git-town sync" in the other worktree

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                       |
      | feature | git fetch --prune --tags      |
      |         | git merge --no-edit --ff main |
      |         | git push                      |
    And the current branch in the other worktree is still "feature"
    And these commits exist now
      | BRANCH  | LOCATION         | MESSAGE                          |
      | main    | origin           | origin main commit               |
      | feature | local            | local feature commit             |
      |         | origin, worktree | origin feature commit            |
      |         |                  | Merge branch 'main' into feature |

  Scenario: undo
    When I run "git-town undo" in the other worktree
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                                |
      | feature | git reset --hard {{ sha-in-worktree-initial 'origin feature commit' }} |
      |         | git push --force-with-lease --force-if-includes                        |
    And the current branch in the other worktree is still "feature"
    And the initial branches and lineage exist now
    And the initial commits exist now
