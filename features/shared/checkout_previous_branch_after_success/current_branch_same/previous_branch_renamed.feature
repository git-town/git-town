Feature: renaming the previous branch makes the main branch the new previous branch

  (see ./previous_branch_same.feature)


  Scenario: rename-branch
    Given I have feature branches named "previous" and "current"
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run `gt rename-branch previous previous-renamed`
    Then I end up on the "current" branch
    And my previous Git branch is now "main"
