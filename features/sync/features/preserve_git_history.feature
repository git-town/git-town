Feature: preserve the previous Git branch

  Scenario:
    Given my repo has the feature branches "previous" and "current"
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run "git-town sync"
    Then I am still on the "current" branch
    And the previous Git branch is still "previous"
