Feature: git extract: handling conflicting remote main branch updates without open changes

  As a developer extracting a commit while there are conflicing changes on the remote main branch
  I want the tool to handle this situation properly
  So that I can use it safely in all edge cases.


  Background:
    Given I am on a feature branch
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT   |
      | main    | remote   | conflicting remote commit | conflicting_file | remote content |
      |         | local    | conflicting local commit  | conflicting_file | local content  |
      | feature | local    | feature commit            | feature_file     |                |
      |         |          | refactor commit           | refactor_file    |                |
    When I run `git extract refactor` with the last commit sha while allowing errors


  Scenario: result
    Then my repo has a rebase in progress
    And there is an abort script for "git extract"


  Scenario: aborting
    When I run `git extract --abort`
    Then I end up on my feature branch
    And there is no "refactor" branch
    And I have the following commits
      | BRANCH  | LOCATION | MESSAGE                   | FILES            |
      | main    | remote   | conflicting remote commit | conflicting_file |
      |         | local    | conflicting local commit  | conflicting_file |
      | feature | local    | feature commit            | feature_file     |
      |         |          | refactor commit           | refactor_file    |
    And there is no rebase in progress
    And there is no abort script for "git extract" anymore


  @finishes-with-non-empty-stash
  Scenario: continuing after resolving the conflicts
    Given TODO: the user should be able to continue here

  @finishes-with-non-empty-stash
  Scenario: continuing without resolving the conflicts
    Given TODO: we should show an error message here and abort
