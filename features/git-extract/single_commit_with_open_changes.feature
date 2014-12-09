Feature: git-extract with a single commit

  As a developer having a feature branch with a commit around an unrelated issue
  I want to be able to extract this commit into its own feature branch
  So that my feature branches remain focussed and code reviews are efficient.


  Background:
    Given I have a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE            | FILE NAME        |
      | main    | remote   | remote main commit | remote_main_file |
      | feature | local    | feature commit     | feature_file     |
      |         |          | refactor commit    | refactor_file    |
    And I am on the "feature" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git extract refactor` with the last commit sha


  Scenario: result
    Then I end up on the "refactor" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And I have the following commits
      | BRANCH   | LOCATION         | MESSAGE            | FILES            |
      | main     | local and remote | remote main commit | remote_main_file |
      | feature  | local            | feature commit     | feature_file     |
      |          |                  | refactor commit    | refactor_file    |
      | refactor | local and remote | remote main commit | remote_main_file |
      |          |                  | refactor commit    | refactor_file    |
    And now I have the following committed files
      | BRANCH   | FILES                           |
      | main     | remote_main_file                |
      | feature  | feature_file, refactor_file     |
      | refactor | remote_main_file, refactor_file |
