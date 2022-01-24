Feature: git town-kill: errors when trying to kill a perennial branch

  (see ../current_branch/on_perennial_branch.feature)

  Background:
    Given my repo has a feature branch named "feature"
    And my repo has the perennial branch "qa"
    And the following commits exist in my repo
      | BRANCH  | LOCATION      | MESSAGE     |
      | feature | local, remote | good commit |
      | qa      | local, remote | qa commit   |
    And I am on the "feature" branch

  Scenario: result
    Given my workspace has an uncommitted file
    When I run "git-town kill qa"
    Then it runs no commands
    And it prints the error:
      """
      you can only kill feature branches
      """
    And I am still on the "feature" branch
    And my workspace still contains my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES          |
      | local      | main, feature, qa |
      | remote     | main, feature, qa |
    And my repo is left with my original commits
