Feature: git extract: don't extract into an already existing branch (with open changes)

  As a developer trying to extract commits into an existing feature branch
  I should get a warning that the target branch already exists
  So that I can extract the commits into a new branch


  Background:
    Given I have feature branches named "feature" and "existing-feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE            | FILE NAME        |
      | main    | remote   | remote main commit | remote_main_file |
      | feature | local    | feature commit     | feature_file     |
      |         |          | refactor commit    | refactor_file    |
    And I am on the "feature" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git extract existing-feature` while allowing errors


  Scenario: result
    Then it runs no Git commands
    And I get the error "A branch named 'existing-feature' already exists"
    And I am still on the "feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
