Feature: git hack: requires a branch name

  As a developer forgetting to provide the name of the new branch to be created
  I should be reminded that I have to provide the branch name to this command
  So that I can use it correctly without having to look that fact up in the readme.


  Background:
    Given I have a feature branch named "existing-feature"
    And the following commits exist in my repository
      | BRANCH           | LOCATION | MESSAGE                 | FILE NAME    |
      | main             | remote   | main commit             | main_file    |
      | existing-feature | local    | existing feature commit | feature_file |
    And I am on the "existing-feature" branch


  Scenario: with open changes
    Given I have an uncommitted file
    When I run `git hack`
    Then it runs no Git commands
    And I get the error "No branch name provided"
    And I am still on the "existing-feature" branch
    And I still have my uncommitted file


  Scenario: without open changes
    When I run `git hack`
    Then it runs no Git commands
    And I get the error "No branch name provided"
    And I am still on the "existing-feature" branch
