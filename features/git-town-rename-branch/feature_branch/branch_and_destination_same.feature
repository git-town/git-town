Feature: git town-rename-branch: does nothing if renaming a feature branch onto itself

  As a developer renaming a feature branch onto itself
  I should get a message saying no action is needed
  So that I am aware that I just did a no-op.


  Background:
    Given my repository has a feature branch named "current-feature"
    And the following commits exist in my repository
      | BRANCH          | LOCATION         | MESSAGE                |
      | current-feature | local and remote | current-feature commit |
    And I am on the "current-feature" branch
    And my workspace has an uncommitted file
    When I run `git-town rename-branch current-feature current-feature`


  Scenario: result
    Then Git Town runs no commands
    And it prints the error "Cannot rename branch to current name."
    And I end up on the "current-feature" branch
    And my workspace still contains my uncommitted file
    And my repository is left with my original commits
