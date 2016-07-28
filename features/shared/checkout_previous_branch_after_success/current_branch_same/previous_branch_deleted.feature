Feature: deleting the current branch makes the main branch the new previous branch

  (see ./previous_branch_same.feature)


  Scenario: git-kill
    Given I have feature branches named "previous" and "current"
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run `git kill previous`
    Then I am still on the "current" branch
    And my previous Git branch is now "main"


  Scenario: git-prune-branches
    Given I have feature branches named "previous" and "current"
    And the "previous" branch gets deleted on the remote
    And the following commit exists in my repository
      | BRANCH  | LOCATION | FILE NAME    | FILE CONTENT    |
      | current | local    | current_file | current content |
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run `git prune-branches`
    Then I am still on the "current" branch
    And my previous Git branch is now "main"


  Scenario: git-ship
    Given I have feature branches named "previous" and "current"
    And the following commit exists in my repository
      | BRANCH   | LOCATION | FILE NAME    | FILE CONTENT    |
      | previous | remote   | feature_file | feature content |
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run `git ship previous -m "feature done"`
    Then I am still on the "current" branch
    And my previous Git branch is now "main"
