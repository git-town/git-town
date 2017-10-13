Feature: git town-kill: errors when trying to kill the main branch

  As a developer accidentally trying to kill the main branch
  I should see an error that I cannot delete the main branch
  So that my main development line remains intact and my project stays shippable.


  Background:
    Given my repository has a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION         | MESSAGE     |
      | feature | local and remote | good commit |
    And I am on the "main" branch


  Scenario: result
    Given my workspace has an uncommitted file
    When I run `git-town kill`
    Then Git Town runs no commands
    And it prints the error "You can only kill feature branches"
    And I am still on the "main" branch
    And my workspace still contains my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES      |
      | local      | main, feature |
      | remote     | main, feature |
    And my repository is left with my original commits
