Feature: git town diff-parent: errors when trying to diff the main branch

  To learn how to use this command correctly
  When accidentally trying to "diff-parent" on the main branch
  I want to see guidance that this isn't possible.

  Scenario:
    Given my repo has a feature branch named "feature"
    And I am on the "main" branch
    When I run "git-town diff-parent"
    Then it runs no commands
    And it prints the error:
      """
      you can only diff-parent feature branches
      """
    And I am still on the "main" branch
