Feature: does not kill perennial branches

  Scenario: trying to delete the main branch
    Given I am on the "main" branch
    When I run "git-town kill"
    Then it runs no commands
    And it prints the error:
      """
      you can only kill feature branches
      """
    And I am still on the "main" branch

  Scenario: trying to delete a perennial branch
    Given my repo has the perennial branch "qa"
    And my repo contains the commits
      | BRANCH | LOCATION      | MESSAGE   |
      | qa     | local, remote | qa commit |
    And I am on the "qa" branch
    When I run "git-town kill"
    Then it runs no commands
    And it prints the error:
      """
      you can only kill feature branches
      """
    And I am still on the "qa" branch
    And the existing branches are
      | REPOSITORY    | BRANCHES |
      | local, remote | main, qa |
