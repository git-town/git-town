Feature: Git checkout history is preserved when the current and previous branch don't change

  As a developer running `git checkout -` after running a Git Town command
  I want to end up on the expected previous branch
  So that Git Town supports my productive use of the Git checkout history


  Scenario: git-kill
    Given I have feature branches named "previous", "current", and "victim"
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run `git kill victim`
    Then I am still on the "current" branch
    And my previous Git branch is still "previous"


  Scenario: git-new-pull-request
    Given I have feature branches named "previous" and "current"
    And I have "open" installed
    And my remote origin is "https://github.com/Originate/git-town.git"
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run `git new-pull-request`
    Then I am still on the "current" branch
    And my previous Git branch is still "previous"


  Scenario: git-prune-branches
    Given I have feature branches named "previous" and "current"
    And the following commit exists in my repository
      | BRANCH   | LOCATION | FILE NAME     | FILE CONTENT     |
      | previous | local    | previous_file | previous content |
      | current  | local    | current_file  | current content  |
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run `git prune-branches`
    Then I am still on the "current" branch
    And my previous Git branch is still "previous"


  Scenario: git-repo
    Given I have feature branches named "previous" and "current"
    And I have "open" installed
    And my remote origin is "https://github.com/Originate/git-town.git"
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run `git repo`
    Then I am still on the "current" branch
    And my previous Git branch is still "previous"


  Scenario: git-ship
    Given I have feature branches named "previous", "current", and "feature"
    And the following commit exists in my repository
      | BRANCH  | LOCATION | FILE NAME    | FILE CONTENT    |
      | feature | remote   | feature_file | feature content |
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run `git ship feature -m "feature done"`
    Then I am still on the "current" branch
    And my previous Git branch is still "previous"


  Scenario: git-sync
    Given I have feature branches named "previous" and "current"
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run `git sync`
    Then I am still on the "current" branch
    And my previous Git branch is still "previous"
