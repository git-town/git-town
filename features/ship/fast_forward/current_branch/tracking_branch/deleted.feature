Feature: shipping a branch whose tracking branch is deleted using the fast-forward strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
    And the current branch is "feature"
    And Git Town setting "ship-strategy" is "fast-forward"
    And origin deletes the "feature" branch
    When I run "git-town ship" and enter "feature done" for the commit message

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git fetch --prune --tags |
    And Git Town prints the error:
      """
      branch "feature" was deleted at the remote
      """
    And the current branch is still "feature"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
