Feature: git extract: refuses to extract into an already existing branch

  As a developer refactoring my branches
  I should not be able to extract commits into already existing feature branches
  So that feature branches remain focussed and code reviews effective


  Background:
    Given I have feature branches named "feature" and "existing"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE            | FILE NAME        |
      | main    | remote   | remote main commit | remote_main_file |
      | feature | local    | feature commit     | feature_file     |
      |         |          | refactor commit    | refactor_file    |
    And I am on the "feature" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git extract existing` while allowing errors


  Scenario: result
    Then it runs no Git commands
    And I get the error "A branch named 'existing' already exists"
    And I am still on the "feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
