Feature: aborting ship of current branch by entering an empty commit message

  Background:
    Given I am on the "feature" branch
    And the following commit exists in my repository
      | BRANCH  | LOCATION | MESSAGE        | FILE NAME    | FILE CONTENT    |
      | feature | local    | feature commit | feature_file | feature content |
    When I run `git ship` and enter an empty commit message


  Scenario: result
    Then I get the error "Aborting ship due to empty commit message"
    And I am still on the "feature" branch
    And I still have the following commits
      | BRANCH  | LOCATION | MESSAGE        | FILES        |
      | feature | local    | feature commit | feature_file |
    And I still have the following committed files
      | BRANCH  | FILES        | CONTENT         |
      | feature | feature_file | feature content |
