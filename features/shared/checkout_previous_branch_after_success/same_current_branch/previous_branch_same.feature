Feature: Allow checking out the correct previous Git branch after running a Git Town commmand that leaves the current and previous branches intact

  As a developer running `git checkout -` after running a Git Town command
  I want to end up on the expected previous branch
  So that I can consistently and effectively use Git's commands


  Scenario: checkout previous branch after git-kill leaves current and previous branches intact
    Given I have feature branches named "previous", "current", and "victim"
    And I am on the "previous" branch
    And I checkout the "current" branch
    And I run `git kill victim`
    When I checkout my previous git branch
    Then I end up on the "previous" branch


  Scenario: checkout previous branch after git-prune-branches leaves current and previous branches intact
    Given I have feature branches named "previous" and "current"
    And the following commit exists in my repository
      | BRANCH   | LOCATION | FILE NAME     | FILE CONTENT     |
      | previous | local    | previous_file | previous content |
      | current  | local    | current_file  | current content  |
    And I am on the "previous" branch
    And I checkout the "current" branch
    And I run `git prune-branches`
    When I checkout my previous git branch
    Then I end up on the "previous" branch


  Scenario: checkout previous branch after git-ship leaves current and previous branches intact
    Given I have feature branches named "previous", "current", and "feature"
    And the following commit exists in my repository
      | BRANCH  | LOCATION | FILE NAME    | FILE CONTENT    |
      | feature | remote   | feature_file | feature content |
    And I am on the "previous" branch
    And I checkout the "current" branch
    And I run `git ship feature -m "feature done"`
    When I checkout my previous git branch
    Then I end up on the "previous" branch


  Scenario: checkout previous branch after git-sync leaves current and previous branches intact
    Given I have feature branches named "previous" and "current"
    And I am on the "previous" branch
    And I checkout the "current" branch
    And I run `git sync`
    When I checkout my previous git branch
    Then I end up on the "previous" branch


  Scenario: checkout previous branch after git-sync-fork leaves current and previous branches intact
    Given my repo has an upstream repo
    And I have feature branches named "previous" and "current"
    And I am on the "previous" branch
    And I checkout the "current" branch
    And I run `git sync-fork`
    When I checkout my previous git branch
    Then I end up on the "previous" branch
