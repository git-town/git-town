Feature: creating a new branch makes the current branch the new previous branch

  (see ../same_current_branch/previous_branch_same.feature)

  Scenario: hack
    Given my repo has the feature branches "previous" and "current"
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run "git-town hack new"
    Then I am now on the "new" branch
    And the previous Git branch is now "current"
