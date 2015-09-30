Feature: Allow checking out the correct previous Git branch after running a Git Town command that deletes the previous and current branches

  (see ../same_current_branch/previous_branch_same.feature)


  Scenario: checkout previous branch after git-prune-branches deletes previous and current branches
    Given I have branches named "previous" and "current"
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run `git prune-branches`
    Then I end up on the "main" branch
    And my previous Git branch is now "main"
