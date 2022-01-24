Feature: deleting the current and previous branches makes the main branch the new previous branch


  Scenario: prune-branches
    Given my repo has the feature branches "previous" and "current"
    And the "previous" branch gets deleted on the remote
    And the "current" branch gets deleted on the remote
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run "git-town prune-branches"
    Then I am now on the "main" branch
    And the previous Git branch is now "main"
