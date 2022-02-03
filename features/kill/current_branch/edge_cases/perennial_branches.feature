Feature: does not kill perennial branches

  Scenario: try delete the main branch
    Given my repo has a feature branch "feature"
    And my repo contains the commits
      | BRANCH  | LOCATION      | MESSAGE     |
      | feature | local, remote | good commit |
    And I am on the "main" branch
    And my workspace has an uncommitted file
    When I run "git-town kill"
    Then it runs no commands
    And it prints the error:
      """
      you can only kill feature branches
      """
    And I am still on the "main" branch
    And my workspace still contains my uncommitted file
    And the existing branches are
      | REPOSITORY    | BRANCHES      |
      | local, remote | main, feature |
    And Git Town still has the original branch hierarchy

  Scenario: try to delete a perennial branch
    Given my repo has a perennial branch "qa"
    And my repo contains the commits
      | BRANCH | LOCATION      | MESSAGE   |
      | qa     | local, remote | qa commit |
    And I am on the "qa" branch
    And my workspace has an uncommitted file
    When I run "git-town kill"
    Then it runs no commands
    And it prints the error:
      """
      you can only kill feature branches
      """
    And I am still on the "qa" branch
    And my workspace still contains my uncommitted file
    And the existing branches are
      | REPOSITORY    | BRANCHES |
      | local, remote | main, qa |
    And Git Town still has no branch hierarchy information
