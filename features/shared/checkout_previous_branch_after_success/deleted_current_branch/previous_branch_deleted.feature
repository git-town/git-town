Feature: Allow checking out the correct previous Git branch after running a Git Town commmand that deletes the previous and current branches

  (see ../same_current_branch/previous_branch_same.feature)


  Scenario: checkout previous branch after a git-prune-branches that deletes previous and current branches
    Given I have feature branches named "previous" and "current"
    And I am on the "previous" branch
    And I checkout the "current" branch
    And I run `git prune-branches`
    When I run `git checkout -` to checkout my previous Git branch
    Then I end up on the "main" branch
