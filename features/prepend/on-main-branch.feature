Feature: git town-prepend: errors when trying to prepend something in front of the main branch

  As a developer accidentally trying to prepend someting in front of the main branch
  I should see an error that the main branch has no parents
  So that I know about my mistake and run "git hack" instead.

  Background:
    Given my repo has a feature branch named "feature"
    And the following commits exist in my repo
      | BRANCH  | LOCATION      | MESSAGE     |
      | feature | local, remote | good commit |
    And I am on the "main" branch

  Scenario: result
    Given my workspace has an uncommitted file
    When I run "git-town prepend new-branch"
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
    And it prints the error:
      """
      the branch "main" is not a feature branch. Only feature branches can have parent branches
      """
    And I am still on the "main" branch
    And my workspace still contains my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES      |
      | local      | main, feature |
      | remote     | main, feature |
    And my repo is left with my original commits
