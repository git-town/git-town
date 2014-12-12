Feature: git ship: allows to resolve remote feature branch conflicts when shipping the current feature branch

  As a developer shipping a feature branch with conflicting remote updates
  I want to get a chance to resolve them
  So that I can ship the branch as planned, and move on to the next feature, and remain productive.


  Background:
    Given I am on the "feature-with-remote-conflicts" branch
    And the following commits exist in my repository
      | BRANCH                        | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT               |
      | feature-with-remote-conflicts | remote   | remote conflicting commit | conflicting_file | remote conflicting content |
      |                               | local    | local conflicting commit  | conflicting_file | local conflicting content  |
    And I run `git ship` while allowing errors


  Scenario: result
    Then I am still on the "feature-with-remote-conflicts" branch
    And my repo has a merge in progress
    And there is an abort script for "git ship"


  Scenario: aborting
    When I run `git ship --abort`
    Then I am still on the "feature-with-remote-conflicts" branch
    And there is no merge in progress
    And there is no abort script for "git ship" anymore
    And I still have the following commits
      | BRANCH                        | LOCATION | MESSAGE                   | FILES            |
      | feature-with-remote-conflicts | local    | local conflicting commit  | conflicting_file |
      |                               | remote   | remote conflicting commit | conflicting_file |
    And I still have the following committed files
      | BRANCH                        | FILES            | CONTENT                   |
      | feature-with-remote-conflicts | conflicting_file | local conflicting content |


  Scenario: continue after resolving the conflict
    Given TODO: make this work


  Scenario: continue without resolving the conflict
    Given TODO: make this work
