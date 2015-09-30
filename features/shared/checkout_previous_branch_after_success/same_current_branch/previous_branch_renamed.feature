Feature: when the previous branch is renamed during a Git Town command, the main branch becomes the new previous branch

  (see ./previous_branch_same.feature)


  Scenario: git-rename-branch
    Given I have branches named "previous" and "current"
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run `git rename-branch previous previous-renamed`
    Then I end up on the "previous-renamed" branch
    And my previous Git branch is now "main"
