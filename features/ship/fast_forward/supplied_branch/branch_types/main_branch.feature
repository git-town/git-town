Feature: does not ship the main branch using the fast-forward strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the current branch is "feature"
    And an uncommitted file
    And Git Town setting "ship-strategy" is "fast-forward"
    When I run "git-town ship main"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git fetch --prune --tags |
    And it prints the error:
      """
      cannot ship the main branch
      """
    And the current branch is still "feature"
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And it prints:
      """
      nothing to undo
      """
    And the current branch is still "feature"
