Feature: Git Ship: resolving conflicts between the supplied feature and main branch (without open changes)

  (see ../current_branch/merge_main_branch_conflict.feature)


  Background:
    Given I have feature branches named "feature" and "other_feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main    | local    | conflicting main commit    | conflicting_file | main content    |
      | feature | local    | conflicting feature commit | conflicting_file | feature content |
    And I am on the "other_feature" branch
    And I run `git ship feature -m "feature done"` while allowing errors


  Scenario: result
    Then it runs the Git commands
      | BRANCH        | COMMAND                            |
      | other_feature | git checkout main                  |
      | main          | git fetch --prune                  |
      | main          | git rebase origin/main             |
      | main          | git push                           |
      | main          | git checkout feature               |
      | feature       | git merge --no-edit origin/feature |
      | feature       | git merge --no-edit main           |
    And I end up on the "feature" branch
    And my repo has a merge in progress


  Scenario: aborting
    When I run `git ship --abort`
    Then it runs the Git commands
      | BRANCH  | COMMAND                    |
      | feature | git merge --abort          |
      | feature | git checkout main          |
      | main    | git checkout other_feature |
    And I end up on the "other_feature" branch
    And there is no merge in progress
    And I still have the following commits
      | BRANCH  | LOCATION         | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main    | local and remote | conflicting main commit    | conflicting_file | main content    |
      | feature | local            | conflicting feature commit | conflicting_file | feature content |


  Scenario: continuing after resolving conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `git ship --continue`
    Then it runs the Git commands
      | BRANCH  | COMMAND                      |
      | feature | git commit --no-edit         |
      | feature | git checkout main            |
      | main    | git merge --squash feature   |
      | main    | git commit -m "feature done" |
      | main    | git push                     |
      | main    | git push origin :feature     |
      | main    | git branch -D feature        |
      | main    | git checkout other_feature   |
    And I end up on the "other_feature" branch
    And there is no "feature" branch
    And I still have the following commits
      | BRANCH | LOCATION         | MESSAGE                 | FILE NAME        |
      | main   | local and remote | feature done            | conflicting_file |
      |        |                  | conflicting main commit | conflicting_file |


  Scenario: continuing after resolving conflicts and comitting
    Given I resolve the conflict in "conflicting_file"
    When I run `git commit --no-edit; git ship --continue`
    Then it runs the Git commands
      | BRANCH  | COMMAND                      |
      | feature | git checkout main            |
      | main    | git merge --squash feature   |
      | main    | git commit -m "feature done" |
      | main    | git push                     |
      | main    | git push origin :feature     |
      | main    | git branch -D feature        |
      | main    | git checkout other_feature   |
    And I end up on the "other_feature" branch
    And there is no "feature" branch
    And I still have the following commits
      | BRANCH | LOCATION         | MESSAGE                 | FILE NAME        |
      | main   | local and remote | feature done            | conflicting_file |
      |        |                  | conflicting main commit | conflicting_file |
