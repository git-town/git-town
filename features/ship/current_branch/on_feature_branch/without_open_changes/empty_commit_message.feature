Feature: git town-ship: aborting the ship of the current feature branch by entering an empty commit message

  As a developer shipping a branch
  I want to be able to abort by entering an empty commit message
  So that shipping has the same experience as committing, and Git Town feels like a natural extension to Git.

  Background:
    Given my repo has a feature branch named "feature"
    And the following commits exist in my repo
      | BRANCH  | LOCATION | MESSAGE        | FILE NAME    | FILE CONTENT    |
      | feature | local    | feature commit | feature_file | feature content |
    And I am on the "feature" branch
    When I run "git-town ship" and enter an empty commit message

  @skipWindows
  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune --tags           |
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
    And it prints the error:
      """
      aborted because commit exited with error
      """
    And I am still on the "feature" branch
    And my repo is left with my original commits

  @skipWindows
  Scenario: undo
    When I run "git-town undo"
    Then it prints the error:
      """
      nothing to undo
      """
    And I am still on the "feature" branch
    And my repo is left with my original commits
