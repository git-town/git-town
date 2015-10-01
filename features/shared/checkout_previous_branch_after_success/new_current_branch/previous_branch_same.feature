Feature: creating a new branch makes the current branch the new previous branch

  (see ../same_current_branch/previous_branch_same.feature)


  Scenario: git-extract
    Given I have branches named "previous" and "current"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE         | FILE NAME     |
      | main    | local    | main commit     | main_file     |
      | current | local    | feature commit  | feature_file  |
      |         |          | refactor commit | refactor_file |
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run `git extract refactor` with the last commit sha
    Then I end up on the "refactor" branch
    And my previous Git branch is now "current"


  Scenario: git-hack
    Given I have branches named "previous" and "current"
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run `git hack new`
    Then I end up on the "new" branch
    And my previous Git branch is now "current"
