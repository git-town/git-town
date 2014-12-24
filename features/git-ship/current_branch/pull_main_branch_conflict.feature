Feature: git ship: resolving conflicts while updating the main branch

  As a developer shipping a branch while there are conflicts between the local and remote main branches
  I want to be given a chance to resolve these conflicts
  So that I can ship the feature as planned and remain productive.


  Background:
    Given I am on the "feature" branch
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT               |
      | main    | remote   | conflicting remote commit | conflicting_file | remote conflicting content |
      |         | local    | conflicting local commit  | conflicting_file | local conflicting content  |
      | feature | local    | feature commit            | feature_file     | feature content            |
    And I run `git ship -m 'feature done'` while allowing errors


  Scenario: result
    Then my repo has a rebase in progress


  Scenario: aborting
    When I run `git ship --abort`
    Then I am still on the "feature" branch
    And there is no rebase in progress
    And I still have the following commits
      | BRANCH  | LOCATION | MESSAGE                   | FILES            |
      | main    | remote   | conflicting remote commit | conflicting_file |
      |         | local    | conflicting local commit  | conflicting_file |
      | feature | local    | feature commit            | feature_file     |
    And I still have the following committed files
      | BRANCH  | FILES            | CONTENT                   |
      | main    | conflicting_file | local conflicting content |
      | feature | feature_file     | feature content           |


  Scenario: continuing after resolving conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `git ship --continue`
    Then it runs the Git commands
      | BRANCH  | COMMAND                            |
      | HEAD    | git rebase --continue              |
      | main    | git push                           |
      | main    | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      | feature | git merge --no-edit main           |
      | feature | git checkout main                  |
      | main    | git merge --squash feature         |
      | main    | git commit -a -m 'feature done'    |
      | main    | git push                           |
      | main    | git push origin :feature           |
      | main    | git branch -D feature              |
    And I end up on the "main" branch
    And there is no "feature" branch
    And I still have the following commits
      | BRANCH  | LOCATION         | MESSAGE                   | FILES            |
      | main    | local and remote | conflicting remote commit | conflicting_file |
      |         |                  | conflicting local commit  | conflicting_file |
      |         |                  | feature done              | feature_file     |
    And now I have the following committed files
      | BRANCH  | FILES                          |
      | main    | conflicting_file, feature_file |


  Scenario: continuing after resolving conflicts and continuing the rebase
    Given I resolve the conflict in "conflicting_file"
    When I run `git rebase --continue; git ship --continue`
    Then it runs the Git commands
      | BRANCH  | COMMAND                            |
      | main    | git push                           |
      | main    | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      | feature | git merge --no-edit main           |
      | feature | git checkout main                  |
      | main    | git merge --squash feature         |
      | main    | git commit -a -m 'feature done'    |
      | main    | git push                           |
      | main    | git push origin :feature           |
      | main    | git branch -D feature              |
    And I end up on the "main" branch
    And there is no "feature" branch
    And I still have the following commits
      | BRANCH  | LOCATION         | MESSAGE                   | FILES            |
      | main    | local and remote | conflicting remote commit | conflicting_file |
      |         |                  | conflicting local commit  | conflicting_file |
      |         |                  | feature done              | feature_file     |
    And now I have the following committed files
      | BRANCH  | FILES                          |
      | main    | conflicting_file, feature_file |
