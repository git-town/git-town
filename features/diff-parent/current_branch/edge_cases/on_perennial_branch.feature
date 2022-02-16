Feature: on perennial branch

  Scenario: on main branch
    And the current branch is "main"
    When I run "git-town diff-parent"
    Then it runs no commands
    And it prints the error:
      """
      you can only diff-parent feature branches
      """

  Scenario: on perennial branch
    And the current branch is a perennial branch "qa"
    When I run "git-town diff-parent"
    Then it runs no commands
    And it prints the error:
      """
      you can only diff-parent feature branches
      """
