Feature: Allow checking out the correct previous Git branch after running a Git Town command that deletes the previous and current branches

  (see ../same_current_branch/previous_branch_same.feature)


  Scenario: checkout previous branch after git-prune-branches deletes previous and current branches
    Given I have branches named "previous" and "current"
    And I am on the "previous" branch
    And I checkout the "current" branch
    When I run `git prune-branches`
    Then I end up on the "main" branch
    When I run `git checkout -`
    Then I end up on the "main" branch
