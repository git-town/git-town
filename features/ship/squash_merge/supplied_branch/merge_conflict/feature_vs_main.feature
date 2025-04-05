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
    And the current branch is "other"
    And Git setting "git-town.ship-strategy" is "squash-merge"
    And I run "git-town ship feature -m 'feature done'"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                         |
      | other  | git fetch --prune --tags        |
      |        | git checkout main               |
      | main   | git merge --squash --ff feature |
      |        | git reset --hard                |
      |        | git checkout other              |
    And Git Town prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    And Git Town prints the error:
      """
      aborted because merge exited with error
      """
    And no merge is in progress

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And Git Town prints:
      """
      nothing to undo
      """
    And no merge is in progress
    And the initial commits exist now
    And the initial branches and lineage exist now
