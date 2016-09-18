Feature: git town-rename-branch: does nothing if renaming a feature branch onto itself

  As a developer renaming a feature branch onto itself
  I should get a message saying no action is needed
  So that I am aware that I just did a no-op.


  Background:
    Given I have a feature branch named "current-feature"
    And the following commits exist in my repository
      | BRANCH          | LOCATION         | MESSAGE                |
      | current-feature | local and remote | current-feature commit |
    And I am on the "current-feature" branch
    And I have an uncommitted file
    When I run `git town-rename-branch current-feature current-feature`


  Scenario: result
    Then I see "Renaming branch to same name, nothing needed."
    And I end up on the "current-feature" branch
    And I still have my uncommitted file
    And I am left with my original commits
