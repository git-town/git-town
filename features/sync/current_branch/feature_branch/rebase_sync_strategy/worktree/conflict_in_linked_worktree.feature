Feature: sync a branch in a "linked worktree" that has a merge conflict

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main    | local, origin | conflicting main commit    | conflicting_file | main content    |
      | feature | local         | conflicting feature commit | conflicting_file | feature content |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the current branch is "main"
    And branch "feature" is active in another worktree
    When I run "git-town sync" in the other worktree

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                         |
      | feature | git fetch --prune --tags                        |
      |         | git push --force-with-lease --force-if-includes |
      |         | git -c rebase.updateRefs=false rebase main      |
    And Git Town prints the error:
      """
      To continue after having resolved conflicts, run "git town continue".
      To go back to where you started, run "git town undo".
      To continue by skipping the current branch, run "git town skip".
      """

  Scenario: undo
    When I run "git-town undo" in the other worktree
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                                       |
      | feature | git rebase --abort                                                            |
      |         | git push --force-with-lease origin {{ sha-initial 'initial commit' }}:feature |
    And the current branch in the other worktree is still "feature"
    And the initial commits exist now

  Scenario: continue with unresolved conflict
    When I run "git-town continue" in the other worktree
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      you must resolve the conflicts before continuing
      """

  Scenario: resolve, commit, and continue
    When I resolve the conflict in "conflicting_file" in the other worktree
    And I run "git rebase --continue" in the other worktree and enter "resolved commit" for the commit message
    And I run "git-town continue" in the other worktree
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                         |
      | feature | git push --force-with-lease --force-if-includes |
    And these commits exist now
      | BRANCH  | LOCATION         | MESSAGE                 | FILE NAME        | FILE CONTENT     |
      | main    | local, origin    | conflicting main commit | conflicting_file | main content     |
      | feature | origin, worktree | resolved commit         | conflicting_file | resolved content |
