Feature: Allow checking out previous git branch to work correctly after running a Git Town commmand

  As a developer running `git checkout -` after running a Git Town command
  I want to end up on the expected previous branch
  So that I can consistently and effectively use git's commands


  Scenario: checkout out previous git branch after git-hack
    Given I have feature branches named "previous" and "current"
    And I am on the "previous" branch
    And I switch to the "current" branch
    And I run `git hack new`
    When I checkout my previous git branch
    Then I end up on the "previous" branch


  Scenario: checkout out previous git branch after git-sync
    Given I have feature branches named "previous" and "current"
    And I am on the "previous" branch
    And I switch to the "current" branch
    And I run `git sync`
    When I checkout my previous git branch
    Then I end up on the "previous" branch


  Scenario: checkout out previous git branch after git-sync-fork
    Given I have feature branches named "previous" and "current"
    And I am on the "previous" branch
    And I switch to the "current" branch
    And I run `git sync`
    When I checkout my previous git branch
    Then I end up on the "previous" branch

