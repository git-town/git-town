Feature: git town-kill: errors when trying to kill the main branch

  (see ../current_branch/on_main_branch.feature)

  Background:
    Given my repo has a feature branch named "feature"
    And the following commits exist in my repo
      | BRANCH  | LOCATION      | MESSAGE     |
      | main    | local, remote | main commit |
      | feature | local, remote | good commit |
    And I am on the "feature" branch

  Scenario: result
    Given my workspace has an uncommitted file
    When I run "git-town kill main"
    Then it runs no commands
    And it prints the error:
      """
      you can only kill feature branches
      """
    And I am still on the "feature" branch
    And my workspace still contains my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES      |
      | local      | main, feature |
      | remote     | main, feature |
    And my repo is left with my original commits
