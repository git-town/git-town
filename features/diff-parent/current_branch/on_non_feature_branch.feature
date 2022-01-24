Feature: git town-diff-parent: errors when trying to diff a perennial branch

  As a developer accidentally trying to diff a perennial branch
  I should see an error that I cannot diff perennial branches
  Because perennial branches cannot have parent branches

  Scenario: result
    Given my repo has the perennial branch "qa"
    And I am on the "qa" branch
    When I run "git-town diff-parent"
    Then it runs no commands
    And it prints the error:
      """
      you can only diff-parent feature branches
      """
    And I am still on the "qa" branch
