Feature: sync a branch in a "linked worktree" that has a merge conflict

  @this
  Scenario:
    Given Git Town setting "sync-feature-strategy" is "rebase"
    And a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE             | FILE NAME        | FILE CONTENT    |
      | main    | local    | local main commit   | conflicting_file | main content    |
      | feature | local    | local parent commit | conflicting_file | feature content |
    And branch "parent" is active in another worktree
    When I run "git-town sync"
    Then it runs the commands
      | BRANCH | COMMAND |
    And the current branch is still "child"
    And these commits exist now
      | BRANCH | LOCATION                | MESSAGE              |
      | main   | local, origin, worktree | origin main commit   |
      |        |                         | local main commit    |
      | child  | local, origin           | origin child commit  |
      |        |                         | origin parent commit |
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
