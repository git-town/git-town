Feature: deleting the current and previous branches makes the main branch the new previous branch

  (see ../same_current_branch/previous_branch_same.feature)


  Scenario: prune-branches
    Given my repository has the feature branches "previous" and "current"
    And the "previous" branch gets deleted on the remote
    And the "current" branch gets deleted on the remote
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run `git-town prune-branches`
    Then I end up on the "main" branch
    And my previous Git branch is now "main"
