Feature: git town-diff-parent: errors when trying to diff a perennial branch

  To learn how to use this command correctly
  When trying to see the changes of a perennial branch
  I should be given guidance that this isn't possible.

  Scenario: result
    And my repo has the perennial branch "qa"
    When I run "git-town diff-parent qa"
    Then it runs no commands
    And it prints the error:
      """
      you can only diff-parent feature branches
      """
