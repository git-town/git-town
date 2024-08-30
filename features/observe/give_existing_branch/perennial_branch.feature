Feature: make another perennial branch an observed branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE      | LOCATIONS     |
      | perennial | perennial | local, origin |
    When I run "git-town observe perennial"

  Scenario: result
    Then it runs no commands
    And it prints the error:
      """
      cannot observe perennial branches
      """
    And the perennial branches are still "perennial"
    And there are still no observed branches

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the perennial branches are still "perennial"
    And there are still no observed branches
