Feature: git hack: resolving conflicting remote main branch updates while moving open changes

  As a developer working on changes that conflict with updates on the main branch
  I want to be given a chance to resolve the differences when moving them to a new feature branch
  So that my work based off the latest state of the code base, I don't run into bigger merge conflicts later, and remain productive.


  Background:
    Given I have a feature branch named "existing_feature"
    Given the following commit exists in my repository
      | BRANCH | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT   |
      | main   | remote   | remote_conflicting_commit | conflicting_file | remote content |
      |        | local    | local_conflicting_commit  | conflicting_file | local content  |
    And I am on the "existing_feature" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git hack new_feature` while allowing errors


  @finishes-with-non-empty-stash
  Scenario: result
    Then it runs the Git commands
      | BRANCH           | COMMAND                |
      | existing_feature | git stash -u           |
      | existing_feature | git checkout main      |
      | main             | git fetch --prune      |
      | main             | git rebase origin/main |
    And my repo has a rebase in progress
    And there is an abort script for "git hack"
    And I don't have an uncommitted file with name: "uncommitted"


  Scenario: aborting
    When I run `git hack --abort`
    Then it runs the Git commands
      | BRANCH           | COMMAND                       |
      | HEAD             | git rebase --abort            |
      | main             | git checkout existing_feature |
      | existing_feature | git stash pop                 |
    And I end up on the "existing_feature" branch
    And I again have an uncommitted file with name: "uncommitted" and content: "stuff"
    And there is no rebase in progress
    And there is no abort script for "git hack" anymore
    And I have the following commits
      | BRANCH | LOCATION | MESSAGE                   | FILES            |
      | main   | remote   | remote_conflicting_commit | conflicting_file |
      |        | local    | local_conflicting_commit  | conflicting_file |
