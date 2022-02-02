Feature: cannot kill perennial branches

  Scenario: trying to delete the main branch
    Given my repo has a feature branch "feature"
    And the following commits exist in my repo
      | BRANCH  | LOCATION      | MESSAGE     |
      | feature | local, remote | good commit |
    And I am on the "main" branch
    Given my workspace has an uncommitted file
    When I run "git-town kill"
    Then it runs no commands
    And it prints the error:
      """
      you can only kill feature branches
      """
    And I am still on the "main" branch
    And my workspace still contains my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES      |
      | local      | main, feature |
      | remote     | main, feature |
    And Git Town is still aware of this branch hierarchy
      | BRANCH  | PARENT |
      | feature | main   |

  Scenario: trying to delete a perennial branch
    Given my repo has the perennial branch "qa"
    And the following commits exist in my repo
      | BRANCH | LOCATION      | MESSAGE   |
      | qa     | local, remote | qa commit |
    And I am on the "qa" branch
    Given my workspace has an uncommitted file
    When I run "git-town kill"
    Then it runs no commands
    And it prints the error:
      """
      you can only kill feature branches
      """
    And I am still on the "qa" branch
    And my workspace still contains my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES |
      | local      | main, qa |
      | remote     | main, qa |
    And Git Town still has no branch hierarchy information
