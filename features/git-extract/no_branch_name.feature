Feature: git extract: errors if no branch name is given

  As a developer forgetting to provide the name of the branch to extract into
  I should see an error explaining the usage of this command
  So that I can use it correctly without having to read the readme again.


  Background:
    Given I have a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE         | FILE NAME     |
      | feature | local    | feature commit  | feature_file  |
      |         |          | refactor commit | refactor_file |
    And I am on the "feature" branch


  Scenario: result
    Given I have an uncommitted file
    When I run `git extract`
    Then it runs no commands
    And I get the error "No branch name provided"
    And I am still on the "feature" branch
    And I still have my uncommitted file

