Feature: git town-rename-branch: does nothing if renaming a feature branch onto itself

  Background:
    Given my repo has a feature branch named "current-feature"
    And the following commits exist in my repo
      | BRANCH          | LOCATION      | MESSAGE                |
      | current-feature | local, remote | current-feature commit |
    And I am on the "current-feature" branch
    And my workspace has an uncommitted file
    When I run "git-town rename-branch current-feature current-feature"

  Scenario: result
    Then it runs no commands
    And it prints the error:
      """
      cannot rename branch to current name
      """
    And I am now on the "current-feature" branch
    And my workspace still contains my uncommitted file
    And my repo is left with my original commits
