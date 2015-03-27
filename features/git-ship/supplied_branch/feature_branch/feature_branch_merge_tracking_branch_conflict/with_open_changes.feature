Feature: git ship: resolving conflicts between the supplied feature branch and its tracking branch (with open changes)

  (see ../../../current_branch/on_feature_branch/without_open_changes/feature_branch_merge_tracking_branch_conflict.feature)


  Background:
    Given I have feature branches named "feature" and "other_feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT               |
      | feature | local    | local conflicting commit  | conflicting_file | local conflicting content  |
      |         | remote   | remote conflicting commit | conflicting_file | remote conflicting content |
    And I am on the "other_feature" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    And I run `git ship feature -m "feature done"`


  @finishes-with-non-empty-stash
  Scenario: result
    Then it runs the Git commands
      | BRANCH        | COMMAND                            |
      | other_feature | git stash -u                       |
      | other_feature | git checkout main                  |
      | main          | git fetch --prune                  |
      | main          | git rebase origin/main             |
      | main          | git checkout feature               |
      | feature       | git merge --no-edit origin/feature |
    And I get the error
      """
      To abort, run "git ship --abort".
      To continue after you have resolved the conflicts, run "git ship --continue".
      """
    And I end up on the "feature" branch
    And I don't have an uncommitted file with name: "uncommitted"
    And my repo has a merge in progress


  Scenario: aborting
    When I run `git ship --abort`
    Then it runs the Git commands
      | BRANCH        | COMMAND                    |
      | feature       | git merge --abort          |
      | feature       | git checkout main          |
      | main          | git checkout other_feature |
      | other_feature | git stash pop              |
    And I end up on the "other_feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And there is no merge in progress
    And I am left with my original commits


  Scenario: continuing after resolving the conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `git ship --continue`
    Then it runs the Git commands
      | BRANCH        | COMMAND                      |
      | feature       | git commit --no-edit         |
      | feature       | git merge --no-edit main     |
      | feature       | git checkout main            |
      | main          | git merge --squash feature   |
      | main          | git commit -m "feature done" |
      | main          | git push                     |
      | main          | git push origin :feature     |
      | main          | git branch -D feature        |
      | main          | git checkout other_feature   |
      | other_feature | git stash pop                |
    And I end up on the "other_feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
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
      | feature       | git checkout main            |
      | main          | git merge --squash feature   |
      | main          | git commit -m "feature done" |
      | main          | git push                     |
      | main          | git push origin :feature     |
      | main          | git branch -D feature        |
      | main          | git checkout other_feature   |
      | other_feature | git stash pop                |
    And I end up on the "other_feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And there is no "feature" branch
    And I still have the following commits
      | BRANCH | LOCATION         | MESSAGE      | FILE NAME        |
      | main   | local and remote | feature done | conflicting_file |
