Feature: when the previous and current branches are deleted during a Git Town command, the main branch becomes the new previous branch

  (see ../same_current_branch/previous_branch_same.feature)


  Scenario: git-prune-branches
    Given I have branches named "previous" and "current"
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run `git prune-branches`
    Then I end up on the "main" branch
    And my previous Git branch is now "main"
