Feature: does not kill perennial branches

  Scenario: main branch
    Given I am on the "main" branch
    When I run "git-town kill"
    Then it runs no commands
    And it prints the error:
      """
      you can only kill feature branches
      """
    And I am still on the "main" branch

  Scenario: perennial branch
    Given my repo has a perennial branch "qa"
    And I am on the "qa" branch
    When I run "git-town kill"
    Then it runs no commands
    And it prints the error:
      """
      you can only kill feature branches
      """
    And I am still on the "qa" branch
