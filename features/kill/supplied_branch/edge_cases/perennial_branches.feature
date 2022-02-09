Feature: does not kill perennial branches

  Scenario: main branch
    Given a feature branch "feature"
    And the current branch is "feature"
    When I run "git-town kill main"
    Then it runs no commands
    And it prints the error:
      """
      you can only kill feature branches
      """
    And the current branch is still "feature"
    And the initial branches and hierarchy exist

  Scenario: perennial branch
    Given a perennial branch "qa"
    And the current branch is "main"
    When I run "git-town kill qa"
    Then it runs no commands
    And it prints the error:
      """
      you can only kill feature branches
      """
