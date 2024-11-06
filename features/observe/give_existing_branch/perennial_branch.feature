Feature: make another perennial branch an observed branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE      | LOCATIONS     |
      | perennial | perennial | local, origin |
    When I run "git-town observe perennial"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      cannot observe perennial branches
      """
    And the perennial branches are still "perennial"
    And there are still no observed branches

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the perennial branches are still "perennial"
    And there are still no observed branches
