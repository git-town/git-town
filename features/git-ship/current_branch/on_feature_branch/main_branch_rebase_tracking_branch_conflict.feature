Feature: git ship: resolving conflicts between the main branch and its tracking branch

  As a developer shipping a branch while there are conflicts between the local and remote main branches
  I want to be given the choice to resolve the conflicts or abort
  So that I can finish the operation as planned or postpone it to a better time.


  Background:
    Given I have a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT               |
      | main    | local    | conflicting local commit  | conflicting_file | local conflicting content  |
      |         | remote   | conflicting remote commit | conflicting_file | remote conflicting content |
      | feature | local    | feature commit            | feature_file     | feature content            |
    And I am on the "feature" branch
    When I run `git ship -m "feature done"`


  Scenario: result
    Then it runs the Git commands
      | BRANCH  | COMMAND                |
      | feature | git fetch --prune      |
      |         | git checkout main      |
      | main    | git rebase origin/main |
    And I get the error
      """
      To abort, run "git ship --abort".
      To continue after you have resolved the conflicts, run "git ship --continue".
      """
    And my repo has a rebase in progress


  Scenario: aborting
    When I run `git ship --abort`
    Then it runs the Git commands
      | BRANCH | COMMAND              |
      | main   | git rebase --abort   |
      |        | git checkout feature |
    And I am still on the "feature" branch
    And there is no rebase in progress
    And I am left with my original commits


  Scenario: continuing after resolving the conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `git ship --continue`
    Then it runs the Git commands
      | BRANCH  | COMMAND                            |
      | main    | git rebase --continue              |
      |         | git push                           |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
      |         | git checkout main                  |
      | main    | git merge --squash feature         |
      |         | git commit -m "feature done"       |
      |         | git push                           |
      |         | git push origin :feature           |
      |         | git branch -D feature              |
    And I end up on the "main" branch
    And there is no "feature" branch
    And I still have the following commits
      | BRANCH | LOCATION         | MESSAGE                   | FILE NAME        |
      | main   | local and remote | conflicting remote commit | conflicting_file |
      |        |                  | conflicting local commit  | conflicting_file |
      |        |                  | feature done              | feature_file     |


  Scenario: continuing after resolving the conflicts and continuing the rebase
    Given I resolve the conflict in "conflicting_file"
    When I run `git rebase --continue; git ship --continue`
    Then it runs the Git commands
      | BRANCH  | COMMAND                            |
      | main    | git push                           |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
      |         | git checkout main                  |
      | main    | git merge --squash feature         |
      |         | git commit -m "feature done"       |
      |         | git push                           |
      |         | git push origin :feature           |
      |         | git branch -D feature              |
    And I end up on the "main" branch
    And there is no "feature" branch
    And I still have the following commits
      | BRANCH | LOCATION         | MESSAGE                   | FILE NAME        |
      | main   | local and remote | conflicting remote commit | conflicting_file |
      |        |                  | conflicting local commit  | conflicting_file |
      |        |                  | feature done              | feature_file     |
