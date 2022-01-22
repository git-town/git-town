Feature: git town diff-parent: errors when trying to diff the main branch

  To learn how to use this command correctly
  When trying to see the changes of the main branch
  I should be given guidance that this isn't possible.

  Scenario:
    When I run "git-town diff-parent main"
    Then it runs no commands
    And it prints the error:
      """
      you can only diff-parent feature branches
      """
