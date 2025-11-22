Feature: sync a branch when main is active in another worktree and has updates

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
    And the current branch is "feature"
    And branch "main" is active in another worktree
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git fetch --prune --tags                |
      |         | git merge --no-edit --ff main           |
      |         | git merge --no-edit --ff origin/feature |
      |         | git push                                |
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE                                                    |
      | main    | origin        | origin main commit                                         |
      |         | worktree      | local main commit                                          |
      | feature | local, origin | local feature commit                                       |
      |         |               | Merge branch 'main' into feature                           |
      |         |               | origin feature commit                                      |
      |         |               | Merge remote-tracking branch 'origin/feature' into feature |
      |         | origin        | local main commit                                          |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                                                |
      | feature | git reset --hard {{ sha 'local feature commit' }}                                      |
      |         | git push --force-with-lease origin {{ sha-in-origin 'origin feature commit' }}:feature |
    And the initial commits exist now
