Feature: git rename-branch: errors if the target branch does not exist

  As a developer mistyping the branch name to rename
  I should get an error that the given branch does not exist
  So that I can rename the correct branch


  Background:
    Given the following commits exist in my repository
      | BRANCH | LOCATION         | MESSAGE     |
      | main   | local and remote | main commit |
    And I am on the "main" branch


  Scenario: with open changes
    When I run `git rename-branch non-existing-feature renamed-branch`
    Given I have an uncommitted file
    Then I get the error "There is no branch named 'non-existing-feature'"
    And I end up on the "main" branch
    And I still have my uncommitted file
    And I am left with my original commits


  Scenario: without open changes
    When I run `git rename-branch non-existing-feature renamed-branch`
    Then I get the error "There is no branch named 'non-existing-feature'"
    And I end up on the "main" branch
    And I am left with my original commits
