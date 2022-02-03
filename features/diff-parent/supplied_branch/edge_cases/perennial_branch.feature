Feature: does not diff perennial branches

  Scenario: main branch
    When I run "git-town diff-parent main"
    Then it runs no commands
    And it prints the error:
      """
      you can only diff-parent feature branches
      """

  Scenario: perennial branch
    Given my repo has the perennial branch "qa"
    When I run "git-town diff-parent qa"
    Then it runs no commands
    And it prints the error:
      """
      you can only diff-parent feature branches
      """
