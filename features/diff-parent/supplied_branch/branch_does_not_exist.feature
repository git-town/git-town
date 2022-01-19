Feature: git town-diff-parent: errors if supplied branch does not exist

  As a developer mistyping the branch name to diff against its parent
  I should get an error that the given branch does not exist
  So that I can diff the correct branch

  Scenario: result
    Given I am on the "main" branch
    When I run "git-town diff-parent non-existing-feature"
    Then it runs no commands
    And it prints the error:
      """
      there is no local branch named "non-existing-feature"
      """
    And I am now on the "main" branch
