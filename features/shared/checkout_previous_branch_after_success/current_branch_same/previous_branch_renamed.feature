Feature: rename the previous branch makes the main branch the new previous branch

  Scenario: rename-branch
    Given my repo has the feature branches "previous" and "current"
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run "git-town rename-branch previous previous-renamed"
    Then I am now on the "current" branch
    And the previous Git branch is now "main"
