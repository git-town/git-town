Feature: git town-diff-parent: errors if supplied branch does not exist

  As a developer mistyping the branch name to diff against its parent
  I should get an error that the given branch does not exist
  So that I can diff the correct branch


  Background:
    Given I am on the "main" branch


  Scenario: result
    When I run "git-town diff-parent non-existing-feature"
    Then it runs no commands
    And it prints the error:
      """
      There is no local branch named 'non-existing-feature'
      """
    And I end up on the "main" branch
