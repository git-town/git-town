Feature: git town-rename-branch: does nothing if renaming a perennial branch onto itself

  As a developer renaming a perennial branch onto itself
  I should get a message saying no action is needed
  So that I am aware that I just did a no-op.


  Background:
    Given I have a perennial branch named "production"
    And the following commits exist in my repository
      | BRANCH     | LOCATION         | MESSAGE           |
      | production | local and remote | production commit |
    And I am on the "production" branch
    And I have an uncommitted file
    When I run `git town-rename-branch production production -f`


  Scenario: result
    Then I see "Renaming branch to same name, nothing to do."
    And I end up on the "production" branch
    And I still have my uncommitted file
    And I am left with my original commits
