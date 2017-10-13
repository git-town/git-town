Feature: git town-ship: resolving conflicts between the current feature branch and the main branch

  As a developer shipping a branch that conflicts with the main branch
  I want to be given the choice to resolve the conflicts or abort
  So that I can finish the operation as planned or postpone it to a better time.


  Background:
    Given my repository has a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main    | local    | conflicting main commit    | conflicting_file | main content    |
      | feature | local    | conflicting feature commit | conflicting_file | feature content |
    And I am on the "feature" branch
    And I run `git-town ship -m "feature done"`


  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune                  |
      |         | git checkout main                  |
      | main    | git rebase origin/main             |
      |         | git push                           |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
    And it prints the error:
      """
      To abort, run "git-town ship --abort".
      To continue after you have resolved the conflicts, run "git-town ship --continue".
      """
    And I am still on the "feature" branch
    And my repo has a merge in progress


  Scenario: aborting
    When I run `git-town ship --abort`
    Then Git Town runs the commands
      | BRANCH  | COMMAND              |
      | feature | git merge --abort    |
      |         | git checkout main    |
      | main    | git checkout feature |
    And I am still on the "feature" branch
    And there is no merge in progress
    And my repository still has the following commits
      | BRANCH  | LOCATION         | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main    | local and remote | conflicting main commit    | conflicting_file | main content    |
      | feature | local            | conflicting feature commit | conflicting_file | feature content |


  Scenario: continuing after resolving the conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `git-town ship --continue`
    Then Git Town runs the commands
      | BRANCH  | COMMAND                      |
      | feature | git commit --no-edit         |
      |         | git checkout main            |
      | main    | git merge --squash feature   |
      |         | git commit -m "feature done" |
      |         | git push                     |
      |         | git push origin :feature     |
      |         | git branch -D feature        |
    And I end up on the "main" branch
    And there is no "feature" branch
    And my repository still has the following commits
      | BRANCH | LOCATION         | MESSAGE                 | FILE NAME        |
      | main   | local and remote | conflicting main commit | conflicting_file |
      |        |                  | feature done            | conflicting_file |


  Scenario: continuing after resolving the conflicts and committing
    Given I resolve the conflict in "conflicting_file"
    When I run `git commit --no-edit; git-town ship --continue`
    Then Git Town runs the commands
      | BRANCH  | COMMAND                      |
      | feature | git checkout main            |
      | main    | git merge --squash feature   |
      |         | git commit -m "feature done" |
      |         | git push                     |
      |         | git push origin :feature     |
      |         | git branch -D feature        |
    And I end up on the "main" branch
    And there is no "feature" branch
    And my repository still has the following commits
      | BRANCH | LOCATION         | MESSAGE                 | FILE NAME        |
      | main   | local and remote | conflicting main commit | conflicting_file |
      |        |                  | feature done            | conflicting_file |
