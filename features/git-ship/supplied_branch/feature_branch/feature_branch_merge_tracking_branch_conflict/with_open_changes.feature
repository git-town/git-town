Feature: git ship: resolving conflicts between the supplied feature branch and its tracking branch (with open changes)

  (see ../../../current_branch/on_feature_branch/without_open_changes/feature_branch_merge_tracking_branch_conflict.feature)


  Background:
    Given I have feature branches named "feature" and "other_feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT               |
      | feature | local    | local conflicting commit  | conflicting_file | local conflicting content  |
      |         | remote   | remote conflicting commit | conflicting_file | remote conflicting content |
    And I am on the "other_feature" branch
    And I have an uncommitted file
    And I run `git ship feature -m "feature done"`


  Scenario: result
    Then it runs the Git commands
      | BRANCH        | COMMAND                            |
      | other_feature | git stash -u                       |
      |               | git checkout main                  |
      | main          | git fetch --prune                  |
      |               | git rebase origin/main             |
      |               | git checkout feature               |
      | feature       | git merge --no-edit origin/feature |
    And I get the error
      """
      To abort, run "git ship --abort".
      To continue after you have resolved the conflicts, run "git ship --continue".
      """
    And I end up on the "feature" branch
    And my uncommitted file is stashed
    And my repo has a merge in progress


  Scenario: aborting
    When I run `git ship --abort`
    Then it runs the Git commands
      | BRANCH        | COMMAND                    |
      | feature       | git merge --abort          |
      |               | git checkout main          |
      | main          | git checkout other_feature |
      | other_feature | git stash pop              |
    And I end up on the "other_feature" branch
    And I still have my uncommitted file
    And there is no merge in progress
    And I am left with my original commits


  Scenario: continuing after resolving the conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `git ship --continue`
    Then it runs the Git commands
      | BRANCH        | COMMAND                      |
      | feature       | git commit --no-edit         |
      |               | git merge --no-edit main     |
      |               | git checkout main            |
      | main          | git merge --squash feature   |
      |               | git commit -m "feature done" |
      |               | git push                     |
      |               | git push origin :feature     |
      |               | git branch -D feature        |
      |               | git checkout other_feature   |
      | other_feature | git stash pop                |
    And I end up on the "other_feature" branch
    And I still have my uncommitted file
    And there is no "feature" branch
    And I still have the following commits
      | BRANCH | LOCATION         | MESSAGE      | FILE NAME        |
      | main   | local and remote | feature done | conflicting_file |


  Scenario: continuing after resolving the conflicts and comitting
    Given I resolve the conflict in "conflicting_file"
    When I run `git commit --no-edit; git ship --continue`
    Then it runs the Git commands
      | BRANCH        | COMMAND                      |
      | feature       | git merge --no-edit main     |
      |               | git checkout main            |
      | main          | git merge --squash feature   |
      |               | git commit -m "feature done" |
      |               | git push                     |
      |               | git push origin :feature     |
      |               | git branch -D feature        |
      |               | git checkout other_feature   |
      | other_feature | git stash pop                |
    And I end up on the "other_feature" branch
    And I still have my uncommitted file
    And there is no "feature" branch
    And I still have the following commits
      | BRANCH | LOCATION         | MESSAGE      | FILE NAME        |
      | main   | local and remote | feature done | conflicting_file |
