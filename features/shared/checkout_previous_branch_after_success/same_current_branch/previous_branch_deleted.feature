Feature: Allow checking out previous git branch to work correctly after running a Git Town commmand that leaves the user on the same branch

# done with tests

  Scenario: checkout previous branch after a git-kill that deletes the previous branch
    Given I have feature branches named "previous" and "current"
    And I am on the "previous" branch
    And I switch to the "current" branch
    And I run `git kill previous`
    When I checkout my previous git branch
    Then I end up on the "main" branch


  Scenario: checkout previous branch after a git-prune-branches that deletes the previous branch
    Given I have feature branches named "previous" and "current"
    And the following commit exists in my repository
      | BRANCH   | LOCATION | FILE NAME     | FILE CONTENT    |
      | current  | local    | current_file  | current content |
    And I am on the "previous" branch
    And I switch to the "current" branch
    And I run `git prune-branches`
    When I checkout my previous git branch
    Then I end up on the "main" branch


  Scenario: checkout previous branch after git-ship that deletes the previous branch
    Given I have feature branches named "previous" and "current"
    And the following commit exists in my repository
      | BRANCH   | LOCATION | FILE NAME    | FILE CONTENT    |
      | previous | remote   | feature_file | feature content |
    And I am on the "previous" branch
    And I switch to the "current" branch
    And I run `git ship previous -m "feature done"`
    When I checkout my previous git branch
    Then I end up on the "main" branch
