Feature: Git Ship: handling conflicting remote feature branch updates when shipping the supplied feature branch without open changes


  Background:
    Given I have feature branches named "feature" and "other_feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT               |
      | feature | remote   | remote conflicting commit | conflicting_file | remote conflicting content |
      | feature | local    | local conflicting commit  | conflicting_file | local conflicting content  |
    And I am on the "other_feature" branch
    And I run `git ship feature -m 'feature done'` while allowing errors


  Scenario: result
    Then I end up on the "feature" branch
    And my repo has a merge in progress
    And there is an abort script for "git ship"


  Scenario: aborting
    When I run `git ship --abort`
    Then I end up on the "other_feature" branch
    And there is no merge in progress
    And there is no abort script for "git ship" anymore
    And I still have the following commits
      | BRANCH  | LOCATION | MESSAGE                   | FILES            |
      | feature | local    | local conflicting commit  | conflicting_file |
      | feature | remote   | remote conflicting commit | conflicting_file |
    And I still have the following committed files
      | BRANCH  | FILES            | CONTENT                   |
      | feature | conflicting_file | local conflicting content |
