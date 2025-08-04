@skipWindows
Feature: handle conflicts between the current feature branch and the main branch (in a local repo)
# TODO: This wrongfully assumes this is a phantom merge conflict,
# and resolves it the wrong way. It should stop and let the user resolve this.

  Background:
    Given a local Git repo
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And I ran "git-town hack feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main    | local    | conflicting main commit    | conflicting_file | main content    |
      | feature | local    | conflicting feature commit | conflicting_file | feature content |
    And the current branch is "feature"
    When I run "git-town sync"

  @debug @this
  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                                      |
      | feature | git -c rebase.updateRefs=false rebase --onto main {{ sha 'initial commit' }} |
    And Git Town prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    And file "conflicting_file" now has content:
      """
      <<<<<<< HEAD
      main content
      =======
      feature content
      >>>>>>> {{ sha-short 'conflicting feature commit' }} (conflicting feature commit)
      """
    And a rebase is now in progress

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                 |
      | feature | git reset --hard {{ sha 'conflicting feature commit' }} |
    And no rebase is now in progress
    And the initial commits exist now
