Feature: Set the main branch as the previous branch after running a Git Town commmand that deletes the initial previous branch

  (see ../same_current_branch/previous_branch_same.feature)


  Scenario: checkout previous branch after a git-prune-branches that deletes previous and current branches
    Given I have feature branches named "previous" and "current"
    And I am on the "previous" branch
    And I switch to the "current" branch
    And I run `git prune-branches`
    When I checkout my previous git branch
    Then I end up on the "main" branch
