Feature: Allow checking out the correct previous Git branch after running a Git Town commmand that creates a new current branch intact and leaves the previous branch intact

  (see ../same_current_branch/previous_branch_same.feature)


  Scenario: checkout previous branch after git-extract
    Given I have feature branches named "previous" and "current"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE         | FILE NAME     |
      | main    | local    | main commit     | main_file     |
      | current | local    | feature commit  | feature_file  |
      |         |          | refactor commit | refactor_file |
    And I am on the "previous" branch
    And I checkout the "current" branch
    When I run `git extract refactor` with the last commit sha
    Then I end up on the "refactor" branch
    When I run `git checkout -`
    Then I end up on the "current" branch


  Scenario: checkout previous branch after git-hack
    Given I have feature branches named "previous" and "current"
    And I am on the "previous" branch
    And I checkout the "current" branch
    When I run `git hack new`
    Then I end up on the "new" branch
    When I run `git checkout -`
    Then I end up on the "current" branch
