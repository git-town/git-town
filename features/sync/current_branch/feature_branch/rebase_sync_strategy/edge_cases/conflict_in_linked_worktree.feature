Feature: sync a branch in a "linked worktree" that has a merge conflict

  Background:
    Given Git Town setting "sync-feature-strategy" is "rebase"
    And a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE             | FILE NAME        | FILE CONTENT    |
      | main    | origin   | local main commit   | conflicting_file | main content    |
      | feature | local    | local parent commit | conflicting_file | feature content |
    And branch "feature" is active in another worktree
    When I run "git-town sync" in the other worktree

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git fetch --prune --tags |
      |         | git rebase origin/main   |
    And it prints the error:
      """
      To continue after having resolved conflicts, run "git town continue".
      To go back to where you started, run "git town undo".
      To continue by skipping the current branch, run "git town skip".
      """

  @debug @this
  Scenario: undo
    When I run "git-town undo" in the other worktree
    Then it runs the commands
      | BRANCH  | COMMAND            |
      | feature | git rebase --abort |
      |         | git stash pop      |
    And the current branch is still "feature"
    And the uncommitted file still exists
    And no rebase is in progress
    And the initial commits exist

  Scenario: continue with unresolved conflict
    When I run "git-town continue"
    Then it runs no commands
    And it prints the error:
      """
      you must resolve the conflicts before continuing
      """
    And the current branch is still "feature"
    And the uncommitted file is stashed
    And a rebase is now in progress




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
