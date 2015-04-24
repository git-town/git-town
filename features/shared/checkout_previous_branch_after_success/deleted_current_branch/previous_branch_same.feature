Feature: Allow checking out previous git branch to work correctly after running a Git Town commmand that deletes the current branch

  (see ../same_current_branch/previous_branch_same.feature)


  Scenario: checkout previous branch after git-kill deletes current branch
    Given I have feature branches named "previous" and "current"
    And I am on the "previous" branch
    And I checkout the "current" branch
    And I run `git kill`
    When I checkout my previous git branch
    Then I end up on the "previous" branch


  Scenario: checkout previous branch after a git-prune-branches deletes current branch
    Given I have feature branches named "previous" and "current"
    And the following commit exists in my repository
      | BRANCH   | LOCATION | FILE NAME     | FILE CONTENT     |
      | previous | local    | previous_file | previous content |
    And I am on the "previous" branch
    And I checkout the "current" branch
    And I run `git prune-branches`
    When I checkout my previous git branch
    Then I end up on the "previous" branch


  Scenario: checkout previous branch after git-ship deletes current branch
    Given I have feature branches named "previous", "current"
    And the following commit exists in my repository
      | BRANCH  | LOCATION | FILE NAME    | FILE CONTENT    |
      | current | remote   | feature_file | feature content |
    And I am on the "previous" branch
    And I checkout the "current" branch
    And I run `git ship -m "feature done"`
    When I checkout my previous git branch
    Then I end up on the "previous" branch
