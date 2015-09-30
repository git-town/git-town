Feature: Git Town commands that don't change the current and previous branch preserve Git history

  As a developer running `git checkout -` after running a Git Town command
  I want to end up on the expected previous branch
  So that Git Town does not interfere with my productive use of the Git command history


  Scenario: git-kill
    Given I have branches named "previous", "current", and "victim"
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run `git kill victim`
    Then I end up on the "current" branch
    And my previous Git branch is now "previous"


  Scenario: git-prune-branches
    Given I have branches named "previous" and "current"
    And the following commit exists in my repository
      | BRANCH   | LOCATION | FILE NAME     | FILE CONTENT     |
      | previous | local    | previous_file | previous content |
      | current  | local    | current_file  | current content  |
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run `git prune-branches`
    Then I end up on the "current" branch
    And my previous Git branch is now "previous"


  Scenario: git-ship
    Given I have branches named "previous", "current", and "feature"
    And the following commit exists in my repository
      | BRANCH  | LOCATION | FILE NAME    | FILE CONTENT    |
      | feature | remote   | feature_file | feature content |
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run `git ship feature -m "feature done"`
    Then I end up on the "current" branch
    And my previous Git branch is now "previous"


  Scenario: git-sync
    Given I have branches named "previous" and "current"
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run `git sync`
    Then I end up on the "current" branch
    And my previous Git branch is now "previous"


  Scenario: git-sync-fork
    Given my repo has an upstream repo
    And I have branches named "previous" and "current"
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run `git sync-fork`
    Then I end up on the "current" branch
    And my previous Git branch is now "previous"
