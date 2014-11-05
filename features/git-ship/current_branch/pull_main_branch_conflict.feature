Feature: git-ship
  the current branch
  conflict when pulling the main branch


  Background:
    Given I am on the "feature" branch
    And the following commits exist in my repository
      | branch  | location | message                   | file name        | file content               |
      | main    | remote   | conflicting remote commit | conflicting_file | remote conflicting content |
      | main    | local    | conflicting local commit  | conflicting_file | local conflicting content  |
      | feature | local    | feature commit            | feature_file     | feature content            |
    And I run `git ship` while allowing errors


  Scenario: result
    Then my repo has a rebase in progress
    And there is an abort script for "git ship"


  Scenario: aborting
    When I run `git ship --abort`
    Then I am still on the "feature" branch
    And there is no rebase in progress
    And there is no abort script for "git ship" anymore
    And I still have the following commits
      | branch  | location | message                   | files            |
      | main    | remote   | conflicting remote commit | conflicting_file |
      | main    | local    | conflicting local commit  | conflicting_file |
      | feature | local    | feature commit            | feature_file     |
    And I still have the following committed files
      | branch  | files            | content                   |
      | main    | conflicting_file | local conflicting content |
      | feature | feature_file     | feature content           |
