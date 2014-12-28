Feature: git extract: abort if no branch name is given (without open changes)

  (see ./no_branch_name_with_open_changes.feature)


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
