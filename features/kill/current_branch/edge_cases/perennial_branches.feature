Feature: does not kill perennial branches

  Scenario: main branch
    Given the current branch is "main"
    When I run "git-town kill"
    Then it runs no commands
    And it prints the error:
      """
      you can only kill feature branches
      """
    And the current branch is still "main"

  Scenario: perennial branch
    Given a perennial branch "qa"
    And the current branch is "qa"
    When I run "git-town kill"
    Then it runs no commands
    And it prints the error:
      """
      you can only kill feature branches
      """
    And the current branch is still "qa"
