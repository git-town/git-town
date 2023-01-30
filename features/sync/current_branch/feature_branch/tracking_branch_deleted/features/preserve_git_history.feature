Feature: preserve the previous Git branch

  Background:
    Given the feature branches "previous" and "current"
    And the current branch is "current" and the previous branch is "previous"

  Scenario: current branch gone, previous branch exists
    And origin deletes the "current" branch
    When I run "git-town sync"
    Then the current branch is now "main"
    And the previous Git branch is still "previous"
