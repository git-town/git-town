Feature: does not kill perennial branches

  Scenario: try to delete the main branch
    Given my repo has a feature branch "feature"
    And my repo contains the commits
      | BRANCH  | LOCATION      | MESSAGE     |
      | main    | local, remote | main commit |
      | feature | local, remote | good commit |
    And I am on the "feature" branch
    And my workspace has an uncommitted file
    When I run "git-town kill main"
    Then it runs no commands
    And it prints the error:
      """
      you can only kill feature branches
      """
    And I am still on the "feature" branch
    And my workspace still contains my uncommitted file
    And the existing branches are
      | REPOSITORY    | BRANCHES      |
      | local, remote | main, feature |
    And my repo is left with my original commits
    And Git Town is now aware of this branch hierarchy
      | BRANCH  | PARENT |
      | feature | main   |

  Scenario: try to delete a perennial branch
    Given my repo has a feature branch "feature"
    And my repo has a perennial branch "qa"
    And my repo contains the commits
      | BRANCH  | LOCATION      | MESSAGE     |
      | feature | local, remote | good commit |
      | qa      | local, remote | qa commit   |
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
    And my repo now has the initial branches
    And my repo is left with my original commits
    And Git Town still has the original branch hierarchy
