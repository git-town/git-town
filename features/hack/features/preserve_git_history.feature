Feature: preserve the previous Git branch

  Scenario:
    Given the feature branches "previous" and "current"
    And the current branch is "current" and the previous branch is "previous"
    When I run "git-town hack new"
    Then the current branch is now "new"
    And the previous Git branch is now "current"
