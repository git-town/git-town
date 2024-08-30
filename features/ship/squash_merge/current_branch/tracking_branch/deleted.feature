@skipWindows
Feature: shipping a branch whose tracking branch is deleted

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
    And Git Town setting "ship-strategy" is "squash-merge"
    And the current branch is "feature"
    And origin deletes the "feature" branch
    When I run "git-town ship" and enter "feature done" for the commit message

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git fetch --prune --tags |
    And it prints the error:
      """
      branch "feature" was deleted at the remote
      """
    And the current branch is still "feature"

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
