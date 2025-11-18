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

  @this
  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                           |
      | child  | git fetch --prune --tags                          |
      |        | git checkout main                                 |
      | main   | git -c rebase.updateRefs=false rebase origin/main |
      |        | git push                                          |
      |        | git checkout child                                |
      | child  | git merge --no-edit --ff origin/parent            |
      |        | git merge --no-edit --ff origin/child             |
      |        | git push                                          |
    And these commits exist now
      | BRANCH | LOCATION                | MESSAGE                                                 |
      | main   | local, origin, worktree | origin main commit                                      |
      |        |                         | local main commit                                       |
      | parent | origin                  | origin parent commit                                    |
      |        | worktree                | local parent commit                                     |
      | child  | local, origin           | local child commit                                      |
      |        | local                   | origin parent commit                                    |
      |        | local, origin           | Merge remote-tracking branch 'origin/parent' into child |
      |        |                         | origin child commit                                     |
      |        |                         | Merge remote-tracking branch 'origin/child' into child  |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                                            |
      | child  | git reset --hard {{ sha 'local child commit' }}                                    |
      |        | git push --force-with-lease origin {{ sha-in-origin 'origin child commit' }}:child |
    And these commits exist now
      | BRANCH | LOCATION                | MESSAGE              |
      | main   | local, origin, worktree | origin main commit   |
      |        |                         | local main commit    |
      | parent | origin                  | origin parent commit |
      |        | worktree                | local parent commit  |
      | child  | local                   | local child commit   |
      |        | origin                  | origin child commit  |
