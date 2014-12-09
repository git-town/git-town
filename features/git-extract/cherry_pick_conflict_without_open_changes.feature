Feature: git-extract lets the user resolve extracts that conflict with changes in main

  As a developer extracting a commit that conflicts with the main branch
  I want to get a chance to resolve these conflicts
  So that I can finish the extract as I originally planned.


  Background:
    Given I have a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE            | FILE NAME        | FILE CONTENT     |
      | main    | local    | conflicting commit | conflicting_file | main content     |
      | feature | local    | feature commit     | feature_file     |                  |
      |         |          | refactor commit    | conflicting_file | refactor content |
    And I am on the "feature" branch
    When I run `git extract refactor` with the last commit sha while allowing errors


  @finishes-with-non-empty-stash
  Scenario: result
    Then I end up on the "refactor" branch
    And my repo has a cherry-pick in progress
    And there is an abort script for "git extract"


  Scenario: aborting
    When I run `git extract --abort`
    Then I end up on the "feature" branch
    And there is no "refactor" branch
    And I have the following commits
      | BRANCH   | LOCATION         | MESSAGE            | FILES            |
      | main     | local and remote | conflicting commit | conflicting_file |
      | feature  | local            | feature commit     | feature_file     |
      |          |                  | refactor commit    | conflicting_file |
    And my repo has no cherry-pick in progress
    And there is no abort script for "git extract" anymore
