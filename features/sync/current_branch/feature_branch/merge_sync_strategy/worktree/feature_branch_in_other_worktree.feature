Feature: Sync a feature branch that is in another worktree than the main branch

  Background:
    Given a feature branch "feature"
    And the perennial branches are "main"
    And the commits
      | BRANCH  | LOCATION | MESSAGE               |
      | main    | local    | local main commit     |
      |         | origin   | origin main commit    |
      | feature | local    | local feature commit  |
      |         | origin   | origin feature commit |
    And branch "feature" is active in another worktree
    When I run "git-town sync" in the other worktree

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune --tags           |
      |         | git merge --no-edit origin/feature |
      |         | git merge --no-edit origin/main    |
      |         | git push                           |
    And the current branch in the other worktree is still "feature"
    And these commits exist now
      | BRANCH  | LOCATION         | MESSAGE                                                    |
      | main    | local            | local main commit                                          |
      |         | origin           | origin main commit                                         |
      | feature | origin, worktree | local feature commit                                       |
      |         |                  | origin feature commit                                      |
      |         |                  | Merge remote-tracking branch 'origin/feature' into feature |
      |         |                  | origin main commit                                         |
      |         |                  | Merge remote-tracking branch 'origin/main' into feature    |

  Scenario: undo
    When I run "git-town undo" in the other worktree
    Then it runs the commands
      | BRANCH  | COMMAND                                                                                |
      | feature | git reset --hard {{ sha 'local feature commit' }}                                      |
      |         | git push --force-with-lease origin {{ sha-in-origin 'origin feature commit' }}:feature |
    And the current branch in the other worktree is still "feature"
    And these commits exist now
      | BRANCH  | LOCATION | MESSAGE               |
      | main    | local    | local main commit     |
      |         | origin   | origin main commit    |
      | feature | origin   | origin feature commit |
      |         | worktree | local feature commit  |
    And the initial branches and lineage exist
