Feature: aborting ship of supplied branch by entering an empty commit message without open changes


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
