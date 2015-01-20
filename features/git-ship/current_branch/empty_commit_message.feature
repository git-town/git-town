Feature: git ship: aborting the ship of the current feature branch by entering an empty commit message

  As a developer shipping a branch
  I want to be able to abort by entering an empty commit message
  So that shipping has the same experience as committing, and Git Town feels like a natural extension to Git.


  Background:
    Given I have a feature branch named "feature"
    And the following commit exists in my repository
      | BRANCH  | LOCATION | MESSAGE        | FILE NAME    | FILE CONTENT    |
      | feature | local    | feature commit | feature_file | feature content |
    And I am on the "feature" branch
    When I run `git ship` and enter an empty commit message


  Scenario: result
    Then it runs the Git commands
      | BRANCH  | COMMAND                            |
      | feature | git checkout main                  |
      | main    | git fetch --prune                  |
      | main    | git rebase origin/main             |
      | main    | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      | feature | git merge --no-edit main           |
      | feature | git checkout main                  |
      | main    | git merge --squash feature         |
      | main    | git commit                         |
      | main    | git reset --hard                   |
      | main    | git checkout feature               |
      | feature | git checkout main                  |
      | main    | git checkout feature               |
    And I get the error "Aborting ship due to empty commit message"
    And I am still on the "feature" branch
    And I still have the following commits
      | BRANCH  | LOCATION | MESSAGE        | FILE NAME    | FILE CONTENT    |
      | feature | local    | feature commit | feature_file | feature content |
