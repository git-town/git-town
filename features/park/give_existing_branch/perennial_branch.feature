Feature: park another perennial branch

  Background:
    Given a Git repo with origin
    And the branch
      | NAME      | TYPE      | LOCATIONS |
      | perennial | perennial | local     |
    When I run "git-town park perennial"

  Scenario: result
    Then it runs no commands
    And it prints the error:
      """
      cannot park perennial branches
      """
    And the perennial branches are still "perennial"
    And there are still no parked branches

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the perennial branches are still "perennial"
    And there are still no observed branches
