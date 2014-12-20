Feature: git extract: don't extract into an already existing branch (without open changes )

  As a developer trying to extract commits into already existing feature branches
  I should get a warning that the target branch already exists
  So that feature branches remain focused and code reviews effective


  Background:
    Given I have feature branches named "feature" and "existing-feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE            | FILE NAME        |
      | main    | remote   | remote main commit | remote_main_file |
      | feature | local    | feature commit     | feature_file     |
      |         |          | refactor commit    | refactor_file    |
    And I am on the "feature" branch
    When I run `git extract existing-feature` while allowing errors


  Scenario: result
    Then it runs no Git commands
    And I get the error "A branch named 'existing-feature' already exists"
    And I am still on the "feature" branch
