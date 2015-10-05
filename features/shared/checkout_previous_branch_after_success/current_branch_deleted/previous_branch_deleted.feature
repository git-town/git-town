Feature: deleting the current and previous branches makes the main branch the new previous branch

  (see ../same_current_branch/previous_branch_same.feature)


  Scenario: git-prune-branches
    Given I have branches named "previous" and "current"
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run `git prune-branches`
    Then I end up on the "main" branch
    And my previous Git branch is now "main"
