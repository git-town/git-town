Feature: Sync a feature branch that is in another worktree than the main branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION | MESSAGE               |
      | main    | local    | local main commit     |
      |         | origin   | origin main commit    |
      | feature | local    | local feature commit  |
      |         | origin   | origin feature commit |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the current branch is "main"
    And branch "feature" is active in another worktree
    When I run "git-town sync" in the other worktree

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                              |
      | feature | git fetch --prune --tags                             |
      |         | git push --force-with-lease --force-if-includes      |
      |         | git -c rebase.updateRefs=false rebase origin/feature |
      |         | git -c rebase.updateRefs=false rebase main           |
      |         | git push --force-with-lease --force-if-includes      |
    And the current branch in the other worktree is still "feature"
    And these commits exist now
      | BRANCH  | LOCATION         | MESSAGE               |
      | main    | local            | local main commit     |
      |         | origin           | origin main commit    |
      | feature | origin           | local main commit     |
      |         | origin, worktree | origin feature commit |
      |         |                  | local feature commit  |

  Scenario: undo
    When I run "git-town undo" in the other worktree
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                                                |
      | feature | git reset --hard {{ sha-in-worktree 'local feature commit' }}                          |
      |         | git push --force-with-lease origin {{ sha-in-origin 'origin feature commit' }}:feature |
    And the current branch in the other worktree is still "feature"
    And the initial branches and lineage exist now
    And the initial commits exist now
