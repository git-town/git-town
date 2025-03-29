Feature: swapping a perennial branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE      | PARENT   | LOCATIONS     |
      | branch-1  | feature   | main     | local, origin |
      | perennial | perennial | branch-1 | local, origin |
    And the current branch is "perennial"
    When I run "git-town swap"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH    | COMMAND                  |
      | perennial | git fetch --prune --tags |
    And Git Town prints the error:
      """
      cannot swap: branch "perennial" is a perennial branch
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the current branch is still "perennial"
    And the initial commits exist now
    And the initial lineage exists now
