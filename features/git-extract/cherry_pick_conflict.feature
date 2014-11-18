Feature: handling cherry-pick conflicts when extracting

  Background:
    Given I am on a feature branch
    And the following commits exist in my repository
      | branch  | location | message            | file name        | file content     |
      | main    | local    | conflicting commit | conflicting_file | main content     |
      | feature | local    | feature commit     | feature_file     |                  |
      | feature | local    | refactor commit    | conflicting_file | refactor content |
    When I run `git extract refactor` with the last commit sha while allowing errors


  Scenario: result
    Then I end up on the "refactor" branch
    And my repo has a cherry-pick in progress
    And there is an abort script for "git extract"


  Scenario: aborting
    When I run `git extract --abort`
    Then I end up on the "feature" branch
    And there is no "refactor" branch
    And I have the following commits
      | branch   | location | message            | files            |
      | main     | local    | conflicting commit | conflicting_file |
      | feature  | local    | feature commit     | feature_file     |
      | feature  | local    | refactor commit    | conflicting_file |
    And my repo has no cherry-pick in progress
    And there is no abort script for "git extract" anymore
