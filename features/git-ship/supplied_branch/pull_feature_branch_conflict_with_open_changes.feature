Feature: git ship: allows to resolve conflicting remote feature branch updates when shipping a given feature branch (with open changes)

  As a developer shipping another feature branch with conflicting remote updates
  I want to get a chance to resolve them
  So that I can ship the branch as planned without further boilerplate Git commands and remain productive by staying focussed on my current feature.


  Background:
    Given I have feature branches named "feature-with-remote-conflicts" and "other_feature"
    And the following commits exist in my repository
      | BRANCH                        | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT               |
      | feature-with-remote-conflicts | remote   | remote conflicting commit | conflicting_file | remote conflicting content |
      |                               | local    | local conflicting commit  | conflicting_file | local conflicting content  |
    And I am on the "other_feature" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    And I run `git ship feature-with-remote-conflicts -m 'feature done'` while allowing errors


  @finishes-with-non-empty-stash
  Scenario: result
    Then I end up on the "feature-with-remote-conflicts" branch
    And I don't have an uncommitted file with name: "uncommitted"
    And my repo has a merge in progress
    And there is an abort script for "git ship"


  Scenario: aborting
    When I run `git ship --abort`
    Then I end up on the "other_feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And there is no merge in progress
    And there is no abort script for "git ship" anymore
    And I still have the following commits
      | BRANCH                        | LOCATION | MESSAGE                   | FILES            |
      | feature-with-remote-conflicts | local    | local conflicting commit  | conflicting_file |
      |                               | remote   | remote conflicting commit | conflicting_file |
    And I still have the following committed files
      | BRANCH                        | FILES            | CONTENT                   |
      | feature-with-remote-conflicts | conflicting_file | local conflicting content |


  @finishes-with-non-empty-stash
  Scenario: continuing
    Given TODO: make this work
