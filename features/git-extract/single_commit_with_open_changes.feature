Feature: git-extract with a single commit

  Background:
    Given I have a feature branch named "feature"
    And the following commits exist in my repository
      | branch  | location | message            | file name        |
      | main    | remote   | remote main commit | remote_main_file |
      | feature | local    | feature commit     | feature_file     |
      | feature | local    | refactor commit    | refactor_file    |
    And I am on the "feature" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git extract refactor` with the last commit sha


  Scenario: result
    Then I end up on the "refactor" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And I have the following commits
      | branch   | location         | message            | files            |
      | main     | local and remote | remote main commit | remote_main_file |
      | feature  | local            | feature commit     | feature_file     |
      | feature  | local            | refactor commit    | refactor_file    |
      | refactor | local and remote | remote main commit | remote_main_file |
      | refactor | local and remote | refactor commit    | refactor_file    |
    And now I have the following committed files
      | branch   | files                           |
      | main     | remote_main_file                |
      | feature  | feature_file, refactor_file     |
      | refactor | remote_main_file, refactor_file |
