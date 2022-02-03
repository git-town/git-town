Feature: does not kill perennial branches

  Scenario: main branch
    Given my repo has a feature branch "feature"
    And I am on the "feature" branch
    When I run "git-town kill main"
    Then it runs no commands
    And it prints the error:
      """
      you can only kill feature branches
      """
    And I am still on the "feature" branch
    And my repo still has its initial branches and branch hierarchy

  Scenario: perennial branch
    Given my repo has a feature branch "feature"
    And my repo has a perennial branch "qa"
    And I am on the "feature" branch
    And my workspace has an uncommitted file
    When I run "git-town kill qa"
    Then it runs no commands
    And it prints the error:
      """
      you can only kill feature branches
      """
    And I am still on the "feature" branch
    And my workspace still contains my uncommitted file
    And my repo still has its initial branches and branch hierarchy
