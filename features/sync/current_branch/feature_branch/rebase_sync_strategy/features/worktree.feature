Feature: sync a branch whose parent is active in another worktree

  Scenario:
    Given Git Town setting "sync-feature-strategy" is "rebase"
    And a feature branch "parent"
    And a feature branch "child" as a child of "parent"
    And the commits
      | BRANCH | LOCATION | MESSAGE              |
      | main   | local    | local main commit    |
      |        | origin   | origin main commit   |
      | parent | local    | local parent commit  |
      |        | origin   | origin parent commit |
      | child  | local    | local child commit   |
      |        | origin   | origin child commit  |
    And branch "parent" is active in another worktree
    And the current branch is "child"
    When I run "git-town sync"
    Then it runs the commands
      | BRANCH | COMMAND                                         |
      | child  | git fetch --prune --tags                        |
      |        | git checkout main                               |
      | main   | git rebase origin/main                          |
      |        | git push                                        |
      |        | git checkout child                              |
      | child  | git rebase origin/child                         |
      |        | git rebase origin/parent                        |
      |        | git push --force-with-lease --force-if-includes |
    And the current branch is still "child"
    And these commits exist now
      | BRANCH | LOCATION                | MESSAGE              |
      | main   | local, origin, worktree | origin main commit   |
      |        |                         | local main commit    |
      | child  | local, origin           | origin parent commit |
      |        |                         | origin child commit  |
      |        |                         | local child commit   |
      | parent | origin                  | origin parent commit |
      |        | worktree                | local parent commit  |

  Scenario: undo
    When I run "git-town undo"
    Then it prints:
      """
      nothing to undo
      """
    And it runs no commands
