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

  Scenario: undo
    When I run "git-town undo -v" in the other worktree
    Then it runs the commands
      | BRANCH  | COMMAND            |
      | feature | git rebase --abort |
      |         | git stash pop      |
    And the current branch is still "feature"
    And the uncommitted file still exists
    And no rebase is in progress
    And the initial commits exist

  Scenario: continue with unresolved conflict
    When I run "git-town continue" in the other worktree
    Then it runs no commands
    And it prints the error:
      """
      you must resolve the conflicts before continuing
      """

  @debug @this
  Scenario: resolve, commit, and continue
    When I resolve the conflict in "conflicting_file" in the other worktree
    And I run "git rebase --continue" in the other worktree and enter "resolved commit" for the commit message
    And I run "git-town continue" in the other worktree
    Then it runs the commands
      | BRANCH  | COMMAND       |
      | feature | git stash pop |
