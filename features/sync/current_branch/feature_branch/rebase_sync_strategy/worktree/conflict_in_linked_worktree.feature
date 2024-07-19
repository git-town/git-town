Feature: sync a branch in a "linked worktree" that has a merge conflict

  Background:
    Given a Git repo clone
    And the branch
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And Git Town setting "sync-feature-strategy" is "rebase"
    And the commits
      | BRANCH  | LOCATION | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main    | origin   | conflicting main commit    | conflicting_file | main content    |
      | feature | local    | conflicting feature commit | conflicting_file | feature content |
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
    When I run "git-town undo" in the other worktree
    Then it runs the commands
      | BRANCH  | COMMAND            |
      | feature | git rebase --abort |
    And the current branch in the other worktree is still "feature"
    And these commits exist now
      | BRANCH  | LOCATION | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main    | origin   | conflicting main commit    | conflicting_file | main content    |
      | feature | worktree | conflicting feature commit | conflicting_file | feature content |

  Scenario: continue with unresolved conflict
    When I run "git-town continue" in the other worktree
    Then it runs no commands
    And it prints the error:
      """
      you must resolve the conflicts before continuing
      """

  Scenario: resolve, commit, and continue
    When I resolve the conflict in "conflicting_file" in the other worktree
    And I run "git rebase --continue" in the other worktree and enter "resolved commit" for the commit message
    And I run "git-town continue" in the other worktree
    Then it runs the commands
      | BRANCH  | COMMAND                                         |
      | feature | git push --force-with-lease --force-if-includes |
    And these commits exist now
      | BRANCH  | LOCATION         | MESSAGE                 | FILE NAME        | FILE CONTENT     |
      | main    | origin           | conflicting main commit | conflicting_file | main content     |
      | feature | origin, worktree | conflicting main commit | conflicting_file | main content     |
      |         |                  | resolved commit         | conflicting_file | resolved content |
