Feature: git ship: allows to resolve merge conflicts between feature and main branch

  As a developer shipping an outdated feature that conflicts with the main branch
  I want to be given a chance to resolve these conflicts
  So that I can ship the feature as I had planned, and remain productive by moving on to the next feature.


  Background:
    Given I am on the "feature" branch
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main    | local    | conflicting main commit    | conflicting_file | main content    |
      | feature | local    | conflicting feature commit | conflicting_file | feature content |
    And I run `git ship` while allowing errors


  Scenario: result
    Then I am still on the "feature" branch
    And my repo has a merge in progress
    And there is an abort script for "git ship"


  Scenario: aborting
    When I run `git ship --abort`
    Then I am still on the "feature" branch
    And there is no merge in progress
    And there is no abort script for "git ship" anymore
    And I still have the following commits
      | BRANCH  | LOCATION         | MESSAGE                    | FILES            |
      | main    | local and remote | conflicting main commit    | conflicting_file |
      | feature | local            | conflicting feature commit | conflicting_file |
    And I still have the following committed files
      | BRANCH  | FILES            | CONTENT         |
      | main    | conflicting_file | main content    |
      | feature | conflicting_file | feature content |


  Scenario: continue after resolving the conflict
    Given TODO: make this work


  Scenario: continue without resolving the conflict
    Given TODO: make this work
