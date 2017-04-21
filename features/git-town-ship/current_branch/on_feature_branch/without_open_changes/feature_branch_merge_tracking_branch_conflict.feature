Feature: git town-ship: resolving conflicts between the current feature branch and its tracking branch

  As a developer shipping a feature branch with conflicting remote updates
  I want to be given the choice to resolve the conflicts or abort
  So that I can finish the operation as planned or postpone it to a better time.


  Background:
    Given I have a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT               |
      | feature | local    | local conflicting commit  | conflicting_file | local conflicting content  |
      |         | remote   | remote conflicting commit | conflicting_file | remote conflicting content |
    And I am on the "feature" branch
    When I run `gt ship -m "feature done"`


  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune                  |
      |         | git checkout main                  |
      | main    | git rebase origin/main             |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
    And I get the error
      """
      To abort, run "gt ship --abort".
      To continue after you have resolved the conflicts, run "gt ship --continue".
      """
    And I am still on the "feature" branch
    And my repo has a merge in progress


  Scenario: aborting
    When I run `gt ship --abort`
    Then it runs the commands
      | BRANCH  | COMMAND              |
      | feature | git merge --abort    |
      |         | git checkout main    |
      | main    | git checkout feature |
    And I am still on the "feature" branch
    And there is no merge in progress
    And I am left with my original commits


  Scenario: continuing after resolving the conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `gt ship --continue`
    Then it runs the commands
      | BRANCH  | COMMAND                      |
      | feature | git commit --no-edit         |
      |         | git merge --no-edit main     |
      |         | git checkout main            |
      | main    | git merge --squash feature   |
      |         | git commit -m "feature done" |
      |         | git push                     |
      |         | git push origin :feature     |
      |         | git branch -D feature        |
    And I end up on the "main" branch
    And there is no "feature" branch
    And I still have the following commits
      | BRANCH | LOCATION         | MESSAGE      | FILE NAME        |
      | main   | local and remote | feature done | conflicting_file |


  Scenario: continuing after resolving the conflicts and committing
    Given I resolve the conflict in "conflicting_file"
    When I run `git commit --no-edit; gt ship --continue`
    Then it runs the commands
      | BRANCH  | COMMAND                      |
      | feature | git merge --no-edit main     |
      |         | git checkout main            |
      | main    | git merge --squash feature   |
      |         | git commit -m "feature done" |
      |         | git push                     |
      |         | git push origin :feature     |
      |         | git branch -D feature        |
    And I end up on the "main" branch
    And there is no "feature" branch
    And I still have the following commits
      | BRANCH | LOCATION         | MESSAGE      | FILE NAME        |
      | main   | local and remote | feature done | conflicting_file |
