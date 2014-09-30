Feature: Git Extract

  Scenario: on the main branch
    Given I am on the main branch
    When I run `git extract` while allowing errors
    Then I am still on the "main" branch


  Scenario: on a feature branch
    Given I am on a feature branch
    And the following commits exist
      | branch  | location | message            | file name        |
      | main    | remote   | remote main commit | remote_main_file |
      | feature | local    | feature commit     | feature_file     |
      | feature | local    | refactor commit    | refactor_file    |
    When I run `git extract refactor` with the last commit sha as an argument
    Then I end up on the "refactor" branch
    And all branches are now synchronized
    And I have the following commits
      | branch   | message            | files            |
      | main     | remote main commit | remote_main_file |
      | feature  | feature commit     | feature_file     |
      | feature  | refactor commit    | refactor_file    |
      | refactor | remote main commit | remote_main_file |
      | refactor | refactor commit    | refactor_file    |
    And now I have the following committed files
      | branch   | name             |
      | main     | remote_main_file |
      | feature  | feature_file     |
      | feature  | refactor_file    |
      | refactor | remote_main_file |
      | refactor | refactor_file    |


  Scenario: user aborts after merge conflict during cherry-picking
    Given I am on a feature branch
    And the following commits exist
      | branch  | location | message            | file name        | file content    |
      | main    | local    | conflicting commit | conflicting_file | main content    |
      | feature | local    | conflicting commit | conflicting_file | feature content |
    When I run `git extract refactor` with the last commit sha as an argument while allowing errors
    Then I end up on the "refactor" branch
    And my repo has a cherry-pick in progress
    And there is an abort script for "git extract"

    When I run `git extract --abort`
    Then I end up on the "feature" branch
    And I have the following commits
      | branch   | message            | files            |
      | main     | conflicting commit | conflicting_file |
      | feature  | conflicting commit | conflicting_file |
    And my repo has no cherry-pick in progress
    And there is no abort script for "git extract" anymore

