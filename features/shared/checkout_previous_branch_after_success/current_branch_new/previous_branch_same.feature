Feature: creating a new branch makes the current branch the new previous branch

  (see ../same_current_branch/previous_branch_same.feature)


  Scenario: git-hack
    Given I have branches named "previous" and "current"
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run `git hack new`
    Then I end up on the "new" branch
    And my previous Git branch is now "current"
