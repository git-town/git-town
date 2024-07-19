Feature: does not diff perennial branches

  Background:
    Given a Git repo clone

  Scenario: main branch
    When I run "git-town diff-parent main"
    Then it runs no commands
    And it prints the error:
      """
      you can only diff-parent feature branches
      """

  Scenario: perennial branch
    Given the branch
      | NAME | TYPE      | LOCATIONS |
      | qa   | perennial | local     |
    When I run "git-town diff-parent qa"
    Then it runs no commands
    And it prints the error:
      """
      you can only diff-parent feature branches
      """
