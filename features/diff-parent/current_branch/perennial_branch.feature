Feature: git town-diff-parent: errors when trying to diff a perennial branch

  To learn how to use this command correctly
  When accidentally trying to "diff-parent" on a perennial branch
  I want to see guidance that this isn't possible.


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
