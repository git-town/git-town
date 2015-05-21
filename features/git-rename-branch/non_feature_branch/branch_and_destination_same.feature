Feature: git rename-branch: does nothing if renaming a non-feature branch onto itself

  As a developer renaming a non-feature branch onto itself
  I should get a message saying no action is needed
  So that I am aware that I just did a no-op.


  Background:
    Given I have a feature branch named "production"
    And my non-feature branches are configured as "production"
    And the following commits exist in my repository
      | BRANCH     | LOCATION         | MESSAGE           |
      | production | local and remote | production commit |
    And I am on the "production" branch


  Scenario: with open changes
    When I run `git rename-branch production production -f`
    Given I have an uncommitted file
    Then I see "Renaming branch to same name, nothing needed."
    And I end up on the "production" branch
    And I still have my uncommitted file
    And I am left with my original commits


  Scenario: without open changes
    When I run `git rename-branch production production -f`
    Then I see "Renaming branch to same name, nothing needed."
    And I end up on the "production" branch
    And I am left with my original commits
