Feature: git ship: aborting the shipping process by entering an empty commit message

  As a developer shipping a branch
  I want to be able to abort by entering an empty commit message
  So that shipping has the same experience as committing, and Git Town feels like a natural extension to Git.


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
