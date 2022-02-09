Feature: preserve the previous Git branch

  Scenario:
    Given the feature branches "previous" and "current"
    And the current branch is "current" and the previous branch is "previous"
    When I run "git-town sync"
    Then the current branch is still "current"
    And the previous Git branch is still "previous"
