Feature: Allow checking out the correct previous Git branch after running a Git Town commmand that leaves the current branch intact and deletes the previous branch

  (see ./previous_branch_same.feature)


  Scenario: checkout previous branch after git-kill leaves current branch intact and deletes the previous branch
    Given I have feature branches named "previous" and "current"
    And I am on the "previous" branch
    And I checkout the "current" branch
    And I run `git kill previous`
    When I run `git checkout -` to checkout my previous Git branch
    Then I end up on the "main" branch


  Scenario: checkout previous branch after git-prune-branches leaves current branch intact and deletes the previous branch
    Given I have feature branches named "previous" and "current"
    And the following commit exists in my repository
      | BRANCH  | LOCATION | FILE NAME    | FILE CONTENT    |
      | current | local    | current_file | current content |
    And I am on the "previous" branch
    And I checkout the "current" branch
    And I run `git prune-branches`
    When I run `git checkout -` to checkout my previous Git branch
    Then I end up on the "main" branch


  Scenario: checkout previous branch after git-ship leaves current branch intact and deletes the previous branch
    Given I have feature branches named "previous" and "current"
    And the following commit exists in my repository
      | BRANCH   | LOCATION | FILE NAME    | FILE CONTENT    |
      | previous | remote   | feature_file | feature content |
    And I am on the "previous" branch
    And I checkout the "current" branch
    And I run `git ship previous -m "feature done"`
    When I run `git checkout -` to checkout my previous Git branch
    Then I end up on the "main" branch
