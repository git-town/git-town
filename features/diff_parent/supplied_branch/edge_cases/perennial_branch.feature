Feature: does not diff perennial branches

  Background:
    Given a Git repo with origin

  Scenario: main branch
    When I run "git-town diff-parent main"
    Then Git Town runs no commands
    And it prints the error:
      """
      you can only diff-parent feature branches
      """

  Scenario: perennial branch
    Given the branches
      | NAME | TYPE      | LOCATIONS |
      | qa   | perennial | local     |
    When I run "git-town diff-parent qa"
    Then Git Town runs no commands
    And it prints the error:
      """
      you can only diff-parent feature branches
      """
