Feature: handle conflicts between the supplied feature branch and the main branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
      | other   | feature | main   | local     |
    And the commits
      | BRANCH  | LOCATION | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main    | local    | conflicting main commit    | conflicting_file | main content    |
      | feature | local    | conflicting feature commit | conflicting_file | feature content |
    And Git setting "git-town.ship-strategy" is "always-merge"
    And the current branch is "other"
    And I run "git-town ship feature" and close the editor

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                             |
      | other  | git fetch --prune --tags            |
      |        | git checkout main                   |
      | main   | git merge --no-ff --edit -- feature |
      |        | git merge --abort                   |
      |        | git checkout other                  |
    And Git Town prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    And Git Town prints the error:
      """
      aborted because merge exited with error
      """
    And no merge is now in progress
  #
  # NOTE: Cannot test undo here.
  # The Git Town command under test has not created an undoable runstate.
  # Executing "git town undo" would undo the Git Town command executed during setup.
