Feature: git ship: abort shipping the given feature branch by entering an empty commit message (without open changes)

  As a user shipping another feature branch
  I want to be able to abort the ship when I realize its the wrong branch by entering an empty commit message
  So that the wrong feature branch doesn't end up on the main development line, my team mates are unaffected by my mistake and remain productive, and my users don't see broken functionality.


  Background:
    Given I have feature branches named "feature" and "other_feature"
    And the following commit exists in my repository
      | BRANCH  | LOCATION | MESSAGE        | FILE NAME    | FILE CONTENT    |
      | feature | local    | feature commit | feature_file | feature content |
    And I am on the "other_feature" branch
    When I run `git ship feature` and enter an empty commit message


  Scenario: result
    Then I get the error "Aborting ship due to empty commit message"
    And I am still on the "other_feature" branch
    And I still have the following commits
      | BRANCH  | LOCATION | MESSAGE        | FILES        |
      | feature | local    | feature commit | feature_file |
    And I still have the following committed files
      | BRANCH  | FILES        | CONTENT         |
      | feature | feature_file | feature content |
