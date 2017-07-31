Feature: git town-ship: aborting the ship of the current feature branch by entering an empty commit message

  As a developer shipping a branch
  I want to be able to abort by entering an empty commit message
  So that shipping has the same experience as committing, and Git Town feels like a natural extension to Git.


  Background:
    Given I have a feature branch named "feature"
    And the following commit exists in my repository
      | BRANCH  | LOCATION | MESSAGE        | FILE NAME    | FILE CONTENT    |
      | feature | local    | feature commit | feature_file | feature content |
    And I am on the "feature" branch
    When I run `git-town ship` and enter an empty commit message


  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune                  |
      |         | git checkout main                  |
      | main    | git rebase origin/main             |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
      |         | git checkout main                  |
      | main    | git merge --squash feature         |
      |         | git commit                         |
      |         | git reset --hard                   |
      |         | git checkout feature               |
      | feature | git checkout main                  |
      | main    | git checkout feature               |
    And I get the error "Aborted because commit exited with error"
    And I am still on the "feature" branch
    And I still have the following commits
      | BRANCH  | LOCATION | MESSAGE        | FILE NAME    | FILE CONTENT    |
      | feature | local    | feature commit | feature_file | feature content |


  Scenario: undo
    When I run `git-town ship --undo`
    Then I get the error "Nothing to undo"
    And I am still on the "feature" branch
    And I still have the following commits
      | BRANCH  | LOCATION | MESSAGE        | FILE NAME    | FILE CONTENT    |
      | feature | local    | feature commit | feature_file | feature content |
