Feature: git-hack handling conflicting remote main branch updates with open changes

  Background:
    Given I have a feature branch named "feature"
    Given the following commit exists in my repository
      | BRANCH | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT   |
      | main   | remote   | remote_conflicting_commit | conflicting_file | remote content |
      |        | local    | local_conflicting_commit  | conflicting_file | local content  |
    And I am on the "feature" branch
    When I run `git hack other_feature` while allowing errors


  Scenario: result
    Then it runs the Git commands
      | BRANCH  | COMMAND                |
      | feature | git checkout main      |
      | main    | git fetch --prune      |
      | main    | git rebase origin/main |
    And my repo has a rebase in progress
    And there is an abort script for "git hack"


  Scenario: aborting
    When I run `git hack --abort`
    Then it runs the Git commands
      | BRANCH  | COMMAND              |
      | HEAD    | git rebase --abort   |
      | main    | git checkout feature |
    And I end up on the "feature" branch
    And there is no rebase in progress
    And there is no abort script for "git hack" anymore
    And I have the following commits
      | BRANCH | LOCATION | MESSAGE                   | FILES            |
      | main   | remote   | remote_conflicting_commit | conflicting_file |
      |        | local    | local_conflicting_commit  | conflicting_file |
