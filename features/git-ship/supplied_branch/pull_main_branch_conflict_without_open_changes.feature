Feature: git ship: resolving conflicting main branch updates when shipping a given feature branch (without open changes)

  (see ../current_branch/pull_main_branch_conflict.feature)


  Background:
    Given I have feature branches named "feature" and "other_feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT               |
      | main    | remote   | conflicting remote commit | conflicting_file | remote conflicting content |
      |         | local    | conflicting local commit  | conflicting_file | local conflicting content  |
      | feature | local    | feature commit            | feature_file     | feature content            |
    And I am on the "other_feature" branch
    And I run `git ship feature -m 'feature done'` while allowing errors


  Scenario: result
    Then my repo has a rebase in progress


  Scenario: aborting
    When I run `git ship --abort`
    Then I am still on the "other_feature" branch
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
      | main    | git checkout other_feature         |
    And I end up on the "other_feature" branch
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
      | main    | git checkout other_feature         |
    And I end up on the "other_feature" branch
    And there is no "feature" branch
    And I still have the following commits
      | BRANCH  | LOCATION         | MESSAGE                   | FILES            |
      | main    | local and remote | conflicting remote commit | conflicting_file |
      |         |                  | conflicting local commit  | conflicting_file |
      |         |                  | feature done              | feature_file     |
    And now I have the following committed files
      | BRANCH  | FILES                          |
      | main    | conflicting_file, feature_file |
