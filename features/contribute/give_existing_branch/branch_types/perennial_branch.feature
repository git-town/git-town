Feature: make another perennial branch a contribution branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE      | LOCATIONS     |
      | perennial | perennial | local, origin |
    When I run "git-town contribute perennial"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      cannot make perennial branches contribution branches
      """
    And branch "perennial" still has type "perennial"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And branch "perennial" still has type "perennial"
