Feature: git town-diff-parent: errors if supplied branch does not exist

  To use the command correctly
  When mistyping the branch name
  I want to get a descriptive error message.

  Scenario:
    When I run "git-town diff-parent non-existing"
    Then it runs no commands
    And it prints the error:
      """
      there is no local branch named "non-existing"
      """
