Feature: git-extract with a single commit

  Background:
    Given I am on a feature branch
    And the following commits exist in my repository
      | branch  | location | message            | file name        |
      | main    | remote   | remote main commit | remote_main_file |
      | feature | local    | feature commit     | feature_file     |
      | feature | local    | refactor commit    | refactor_file    |
    When I run `git extract refactor` with the last commit sha


  Scenario: result
    Then I end up on the "refactor" branch
    And all branches are now synchronized
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
