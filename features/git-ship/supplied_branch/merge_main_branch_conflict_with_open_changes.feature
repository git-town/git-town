Feature: handling merge conflicts between feature and main branch when shipping the supplied feature branch with open changes


  Background:
    Given I have feature branches named "feature" and "other_feature"
    And the following commits exist in my repository
      | branch  | location | message                    | file name        | file content    |
      | main    | local    | conflicting main commit    | conflicting_file | main content    |
      | feature | local    | conflicting feature commit | conflicting_file | feature content |
    And I am on the "other_feature" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    And I run `git ship feature -m 'feature done'` while allowing errors


  @finishes-with-non-empty-stash
  Scenario: result
    Then I end up on the "feature" branch
    And my repo has a merge in progress
    And there is an abort script for "git ship"
    And I don't have an uncommitted file with name: "uncommitted"


  Scenario: aborting
    When I run `git ship --abort`
    Then I end up on the "other_feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And there is no merge in progress
    And there is no abort script for "git ship" anymore
    And I still have the following commits
      | branch  | location | message                    | files            |
      | main    | local    | conflicting main commit    | conflicting_file |
      | feature | local    | conflicting feature commit | conflicting_file |
    And I still have the following committed files
      | branch  | files            | content         |
      | main    | conflicting_file | main content    |
      | feature | conflicting_file | feature content |
