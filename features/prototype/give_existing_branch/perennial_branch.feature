Feature: prototype another perennial branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE      | LOCATIONS |
      | perennial | perennial | local     |
    When I run "git-town prototype perennial"

  Scenario: result
    Then it runs no commands
    And it prints the error:
      """
      cannot prototype perennial branches
      """
    And the perennial branches are still "perennial"
    And there are still no prototype branches

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the perennial branches are still "perennial"
    And there are still no prototype branches
