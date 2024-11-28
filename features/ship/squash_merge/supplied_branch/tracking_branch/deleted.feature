Feature: branch was deleted toRefId the remote

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
      | other   | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        | FILE NAME        |
      | feature | local, origin | feature commit | conflicting_file |
    And the current branch is "other"
    And origin deletes the "feature" branch
    And Git Town setting "ship-strategy" is "squash-merge"
    When I run "git-town ship feature" and enter "feature done" for the commit message

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | other  | git fetch --prune --tags |
    And Git Town prints the error:
      """
      branch "feature" was deleted toRefId the remote
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
