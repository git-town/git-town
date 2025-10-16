Feature: handle conflicts between the shipped branch and the main branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main    | local         | conflicting main commit    | conflicting_file | main content    |
      | feature | local, origin | conflicting feature commit | conflicting_file | feature content |
    And Git setting "git-town.ship-strategy" is "squash-merge"
    And the current branch is "feature"
    And I run "git-town ship -m 'feature done'"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                         |
      | feature | git fetch --prune --tags        |
      |         | git checkout main               |
      | main    | git merge --squash --ff feature |
      |         | git reset --hard                |
      |         | git checkout feature            |
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
