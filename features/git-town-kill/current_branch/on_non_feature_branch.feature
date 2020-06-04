Feature: git town-kill: errors when trying to kill a perennial branch

  As a developer accidentally trying to kill a perennial branch
  I should see an error that I cannot delete perennial branches
  So that my release infrastructure remains intact and my project stays shippable.


  Background:
    Given my repo has the perennial branch "qa"
    And the following commits exist in my repo
      | BRANCH | LOCATION      | MESSAGE   |
      | qa     | local, remote | qa commit |
    And I am on the "qa" branch


  Scenario: result
    Given my workspace has an uncommitted file
    When I run "git-town kill"
    Then it runs no commands
    And it prints the error:
      """
      you can only kill feature branches
      """
    And I am still on the "qa" branch
    And my workspace still contains my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES |
      | local      | main, qa |
      | remote     | main, qa |
    And my repo is left with my original commits
