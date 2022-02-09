Feature: preserve the previous Git branch

  Background:
    Given my repo has the feature branches "previous" and "current"
    And I am on the "current" branch with "previous" as the previous Git branch

  Scenario: current branch renamed
    When I run "git-town rename-branch current new"
    Then I am now on the "new" branch
    And the previous Git branch is still "previous"

  Scenario: previous branch renamed
    When I run "git-town rename-branch previous new"
    Then I am now on the "current" branch
    And the previous Git branch is now "main"
