Feature: git hack: enforces being given a branch name when starting a new feature

  As a developer trying to create a new feature branch but forgetting to provide the new branch name
  I should be reminded about the correct syntax for this command
  So that I can use it correctly without having to look that up the readme.


  Background:
    Given I have a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE        | FILE NAME    |
      | main    | remote   | main commit    | main_file    |
      | feature | local    | feature commit | feature_file |
    And I am on the "feature" branch
    When I run `git hack` while allowing errors


  Scenario: result
    Then it runs no Git commands
    And I get the error "No branch name provided"
    And I am still on the "feature" branch
