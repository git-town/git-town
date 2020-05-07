Feature: git town-diff-parent: errors when trying to diff the main branch

  As a developer accidentally trying to diff the main branch
  I should see an error that I cannot diff the main branch
  Because the master branch cannot have a parent branch


  Background:
    Given my repository has a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION      | MESSAGE     |
      | feature | local, remote | good commit |
    And I am on the "main" branch


  Scenario: result
    Given my workspace has an uncommitted file
    When I run "git-town diff-parent"
    Then it runs no commands
    And it prints the error:
      """
      You can only diff-parent feature branches
      """
    And I am still on the "main" branch
    And my workspace still contains my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES      |
      | local      | main, feature |
      | remote     | main, feature |
    And my repository is left with my original commits
