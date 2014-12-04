Feature: git-hack handling conflicting remote main branch updates with open changes

  Background:
    Given I have a feature branch named "feature"
    Given the following commit exists in my repository
      | BRANCH | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT   |
      | main   | remote   | remote_conflicting_commit | conflicting_file | remote content |
      | main   | local    | local_conflicting_commit  | conflicting_file | local content  |
    And I am on the "feature" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git hack other_feature` while allowing errors


  @finishes-with-non-empty-stash
  Scenario: result
    Then my repo has a rebase in progress
    And there is an abort script for "git hack"
    And I don't have an uncommitted file with name: "uncommitted"


  Scenario: aborting
    When I run `git hack --abort`
    Then I end up on the "feature" branch
    And I again have an uncommitted file with name: "uncommitted" and content: "stuff"
    And there is no rebase in progress
    And there is no abort script for "git hack" anymore
    And I have the following commits
      | BRANCH | LOCATION | MESSAGE                   | FILES            |
      | main   | remote   | remote_conflicting_commit | conflicting_file |
      | main   | local    | local_conflicting_commit  | conflicting_file |
