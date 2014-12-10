Feature: git-hack handles conflicting remote main branch updates while starting a new feature

  As a developer creating a feature branch while there are conflicting changes on the remote main branch
  I want the tool to handle this situation properly
  So that I can use it safely in all edge cases.


  Background:
    Given I have a feature branch named "existing_feature"
    Given the following commit exists in my repository
      | BRANCH | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT   |
      | main   | remote   | remote_conflicting_commit | conflicting_file | remote content |
      |        | local    | local_conflicting_commit  | conflicting_file | local content  |
    And I am on the "existing_feature" branch
    When I run `git hack new_feature` while allowing errors


  Scenario: result
    Then it runs the Git commands
      | BRANCH           | COMMAND                |
      | existing_feature | git checkout main      |
      | main             | git fetch --prune      |
      | main             | git rebase origin/main |
    And my repo has a rebase in progress
    And there is an abort script for "git hack"


  Scenario: aborting
    When I run `git hack --abort`
    Then it runs the Git commands
      | BRANCH  | COMMAND                       |
      | HEAD    | git rebase --abort            |
      | main    | git checkout existing_feature |
    And I end up on the "existing_feature" branch
    And there is no rebase in progress
    And there is no abort script for "git hack" anymore
    And I have the following commits
      | BRANCH | LOCATION | MESSAGE                   | FILES            |
      | main   | remote   | remote_conflicting_commit | conflicting_file |
      |        | local    | local_conflicting_commit  | conflicting_file |
