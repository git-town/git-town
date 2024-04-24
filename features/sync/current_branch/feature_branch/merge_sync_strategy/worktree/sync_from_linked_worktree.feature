Feature: sync a branch whose parent is active in another worktree

  Background:
    Given a feature branch "parent"
    And a feature branch "child" as a child of "parent"
    And the commits
      | BRANCH | LOCATION | MESSAGE              |
      | main   | local    | local main commit    |
      |        | origin   | origin main commit   |
      | parent | local    | local parent commit  |
      |        | origin   | origin parent commit |
      | child  | local    | local child commit   |
      |        | origin   | origin child commit  |
    And branch "child" is active in another worktree
    And the current branch is "parent"
    When I run "git-town sync" in the other worktree

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                                |
      | child  | git fetch --prune --tags               |
      |        | git checkout main                      |
      | main   | git rebase origin/main                 |
      |        | git push                               |
      |        | git checkout child                     |
      | child  | git merge --no-edit --ff origin/child  |
      |        | git merge --no-edit --ff origin/parent |
      |        | git push                               |
    And the current branch is still "parent"
    And these commits exist now
      | BRANCH | LOCATION                | MESSAGE                                                 |
      | main   | local, origin, worktree | origin main commit                                      |
      |        |                         | local main commit                                       |
      | child  | origin, worktree        | local child commit                                      |
      |        |                         | origin child commit                                     |
      |        |                         | Merge remote-tracking branch 'origin/child' into child  |
      |        |                         | origin parent commit                                    |
      |        |                         | Merge remote-tracking branch 'origin/parent' into child |
      | parent | local                   | local parent commit                                     |
      |        | origin                  | origin parent commit                                    |

  Scenario: undo
    When I run "git-town undo" in the other worktree
    Then it runs the commands
      | BRANCH | COMMAND                                                                            |
      | child  | git reset --hard {{ sha 'local child commit' }}                                    |
      |        | git push --force-with-lease origin {{ sha-in-origin 'origin child commit' }}:child |
    And the current branch is still "parent"
    And the current branch in the other worktree is still "child"
    And these commits exist now
      | BRANCH | LOCATION                | MESSAGE              |
      | main   | local, origin, worktree | origin main commit   |
      |        |                         | local main commit    |
      | child  | origin                  | origin child commit  |
      |        | worktree                | local child commit   |
      | parent | local                   | local parent commit  |
      |        | origin                  | origin parent commit |
