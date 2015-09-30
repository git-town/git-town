Feature: Allow checking out the correct previous Git branch after running a Git Town command that keeps the previous branch intact and deletes the current branch

  (see ../same_current_branch/previous_branch_same.feature)


  Scenario: checkout previous branch after git-kill keeps the previous branch intact and deletes the current branch
    Given I have branches named "previous" and "current"
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run `git kill`
    Then I end up on the "main" branch
    And my previous Git branch is now "previous"


  Scenario: checkout previous branch after git-prune-branches keeps the previous branch intact and deletes the current branch
    Given I have branches named "previous" and "current"
    And the following commit exists in my repository
      | BRANCH   | LOCATION | FILE NAME     | FILE CONTENT     |
      | previous | local    | previous_file | previous content |
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run `git prune-branches`
    Then I end up on the "main" branch
    And my previous Git branch is now "previous"


  Scenario: checkout previous branch after git-ship keeps the previous branch intact and deletes the current branch
    Given I have branches named "previous" and "current"
    And the following commit exists in my repository
      | BRANCH  | LOCATION | FILE NAME    | FILE CONTENT    |
      | current | remote   | feature_file | feature content |
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run `git ship -m "feature done"`
    Then I end up on the "main" branch
    And my previous Git branch is now "previous"
