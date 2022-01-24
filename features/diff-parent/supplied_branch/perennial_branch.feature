Feature: git town-diff-parent: errors when trying to diff a perennial branch

  Scenario: result
    Given my repo has a feature branch named "feature"
    And my repo has the perennial branch "qa"
    And I am on the "feature" branch
    When I run "git-town diff-parent qa"
    Then it runs no commands
    And it prints the error:
      """
      you can only diff-parent feature branches
      """
    And I am still on the "feature" branch
