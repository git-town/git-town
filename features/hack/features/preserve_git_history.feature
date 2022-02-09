Feature: preserve the previous Git branch

  Scenario:
    Given the feature branches "previous" and "current"
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run "git-town hack new"
    Then I am now on the "new" branch
    And the previous Git branch is now "current"
