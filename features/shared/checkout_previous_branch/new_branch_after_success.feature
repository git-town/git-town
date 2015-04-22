Feature: Allow checking out previous git branch to work correctly after running a Git Town commmand that creates a new branch

  (see ./same_branch_after_success.feature)


  Scenario: checkout previous git branch after git-extract
    Given I have feature branches named "previous" and "current"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE         | FILE NAME     |
      | main    | local    | main commit     | main_file     |
      | current | local    | feature commit  | feature_file  |
      |         |          | refactor commit | refactor_file |
    And I am on the "previous" branch
    And I switch to the "current" branch
    And I run `git extract refactor` with the last commit sha
    When I checkout my previous git branch
    Then I end up on the "current" branch


  Scenario: checkout previous git branch after git-hack
    Given I have feature branches named "previous" and "current"
    And I am on the "previous" branch
    And I switch to the "current" branch
    And I run `git hack new`
    When I checkout my previous git branch
    Then I end up on the "current" branch
