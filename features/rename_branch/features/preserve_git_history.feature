Feature: preserve the previous Git branch

  Background:
    Given the feature branches "previous" and "current"
    And the current branch is "current" and the previous branch is "previous"

  Scenario: current branch renamed
    When I run "git-town rename-branch current new"
    Then the current branch is now "new"
    And the previous Git branch is still "previous"

  Scenario: previous branch renamed
    When I run "git-town rename-branch previous new"
    Then the current branch is now "current"
    And the previous Git branch is now "new"
