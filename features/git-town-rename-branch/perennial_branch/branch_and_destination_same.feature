Feature: git town-rename-branch: does nothing if renaming a perennial branch onto itself

  As a developer renaming a perennial branch onto itself
  I should get a message saying no action is needed
  So that I am aware that I just did a no-op.


  Background:
    Given my repository has a perennial branch named "production"
    And the following commits exist in my repository
      | BRANCH     | LOCATION         | MESSAGE           |
      | production | local and remote | production commit |
    And I am on the "production" branch
    And my workspace has an uncommitted file
    When I run `git-town rename-branch --force production production`


  Scenario: result
    Then Git Town runs no commands
    And it prints the error "Cannot rename branch to current name."
    And I end up on the "production" branch
    And my workspace still contains my uncommitted file
    And my repository is left with my original commits
