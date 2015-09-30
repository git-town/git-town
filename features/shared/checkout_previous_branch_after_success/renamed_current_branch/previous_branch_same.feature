Feature: when the current branch is renamed during a Git Town command, the previous branch history is preserved

  (see ../same_current_branch/previous_branch_same.feature)


  Scenario: git-rename-branch
    Given I have branches named "previous" and "current"
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run `git rename-branch current current-new`
    Then I end up on the "current-new" branch
    And my previous Git branch is still "previous"
