Feature: git extract errors without a branch name without open changes

  As a developer forgetting to provide the target branch to extract into
  I want to be reminded about the correct syntax for this command
  So that I can use Git Town correctly without having to memorize the syntax.


  Background:
    Given I have a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE            | FILE NAME        |
      | main    | remote   | remote main commit | remote_main_file |
      | feature | local    | feature commit     | feature_file     |
      |         |          | refactor commit    | refactor_file    |
    And I am on the "feature" branch
    When I run `git extract` while allowing errors


  Scenario: result
    Then it runs no Git commands
    And I get the error "No branch name provided"
    And I am still on the "feature" branch
