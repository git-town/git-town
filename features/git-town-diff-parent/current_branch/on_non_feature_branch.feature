Feature: git town-diff-parent: errors when trying to diff a perennial branch

  As a developer accidentally trying to diff a perennial branch
  I should see an error that I cannot diff perennial branches
  Because perennial branches cannot have parent branches


  Background:
    Given my repository has the perennial branch "qa"
    And the following commits exist in my repository
      | BRANCH | LOCATION      | MESSAGE   |
      | qa     | local, remote | qa commit |
    And I am on the "qa" branch


  Scenario: result
    Given my workspace has an uncommitted file
    When I run "git-town diff-parent"
    Then it runs no commands
    And it prints the error:
      """
      You can only diff-parent feature branches
      """
    And I am still on the "qa" branch
    And my workspace still contains my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES |
      | local      | main, qa |
      | remote     | main, qa |
    And my repository is left with my original commits
