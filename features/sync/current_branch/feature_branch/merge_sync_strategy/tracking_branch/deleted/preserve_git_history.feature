Feature: preserve the previous Git branch

  Background:
    Given the feature branches "current", "previous", and "other"
    And the current branch is "current" and the previous branch is "previous"

  Scenario: current branch gone, previous branch exists
    Given origin deletes the "current" branch
    When I run "git-town sync"
    Then the current branch is now "previous"
    And the previous Git branch is still "previous"
