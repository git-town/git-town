Feature: Allow checking out previous git branch to work correctly after running a Git Town commmand that leaves the user on the same branch

  As a developer running `git checkout -` after running a Git Town command
  I want to end up on the expected previous branch
  So that I can consistently and effectively use git's commands


  Scenario: checkout previous git branch after git-kill
    Given I have feature branches named "previous", "current", and "victim"
    And I am on the "previous" branch
    And I switch to the "current" branch
    And I run `git kill victim`
    When I checkout my previous git branch
    Then I end up on the "previous" branch


  Scenario: checkout previous git branch after git-ship
    Given I have feature branches named "previous", "current", and "feature"
    And the following commit exists in my repository
      | BRANCH  | LOCATION | FILE NAME    | FILE CONTENT    |
      | feature | remote   | feature_file | feature content |
    And I am on the "previous" branch
    And I switch to the "current" branch
    And I run `git ship feature -m "feature done"`
    When I checkout my previous git branch
    Then I end up on the "previous" branch


  Scenario: checkout previous git branch after git-sync
    Given I have feature branches named "previous" and "current"
    And I am on the "previous" branch
    And I switch to the "current" branch
    And I run `git sync`
    When I checkout my previous git branch
    Then I end up on the "previous" branch


  Scenario: checkout previous git branch after git-sync-fork
    Given my repo has an upstream repo
    And I have feature branches named "previous" and "current"
    And I am on the "previous" branch
    And I switch to the "current" branch
    And I run `git sync-fork`
    When I checkout my previous git branch
    Then I end up on the "previous" branch

