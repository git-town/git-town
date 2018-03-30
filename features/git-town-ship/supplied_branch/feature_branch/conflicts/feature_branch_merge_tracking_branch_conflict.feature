Feature: git town-ship: resolving conflicts between the supplied feature branch and its tracking branch

  (see ../../../current_branch/on_feature_branch/without_open_changes/feature_branch_merge_tracking_branch_conflict.feature)


  Background:
    Given my repository has the feature branches "feature" and "other-feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT               |
      | feature | local    | local conflicting commit  | conflicting_file | local conflicting content  |
      |         | remote   | remote conflicting commit | conflicting_file | remote conflicting content |
    And I am on the "other-feature" branch
    And my workspace has an uncommitted file
    And I run `git-town ship feature -m "feature done"`


  Scenario: result
    Then it runs the commands
      | BRANCH        | COMMAND                            |
      | other-feature | git fetch --prune                  |
      |               | git add -A                         |
      |               | git stash                          |
      |               | git checkout main                  |
      | main          | git rebase origin/main             |
      |               | git checkout feature               |
      | feature       | git merge --no-edit origin/feature |
    And it prints the error:
      """
      To abort, run "git-town abort".
      To continue after having resolved conflicts, run "git-town continue".
      """
    And I end up on the "feature" branch
    And my uncommitted file is stashed
    And my repo has a merge in progress


  Scenario: aborting
    When I run `git-town abort`
    Then it runs the commands
      | BRANCH        | COMMAND                    |
      | feature       | git merge --abort          |
      |               | git checkout main          |
      | main          | git checkout other-feature |
      | other-feature | git stash pop              |
    And I end up on the "other-feature" branch
    And my workspace still contains my uncommitted file
    And there is no merge in progress
    And my repository is left with my original commits


  Scenario: continuing after resolving the conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `git-town continue`
    Then it runs the commands
      | BRANCH        | COMMAND                      |
      | feature       | git commit --no-edit         |
      |               | git merge --no-edit main     |
      |               | git checkout main            |
      | main          | git merge --squash feature   |
      |               | git commit -m "feature done" |
      |               | git push                     |
      |               | git push origin :feature     |
      |               | git branch -D feature        |
      |               | git checkout other-feature   |
      | other-feature | git stash pop                |
    And I end up on the "other-feature" branch
    And my workspace still contains my uncommitted file
    And there is no "feature" branch
    And my repository still has the following commits
      | BRANCH | LOCATION         | MESSAGE      | FILE NAME        |
      | main   | local and remote | feature done | conflicting_file |


  Scenario: continuing after resolving the conflicts and comitting
    Given I resolve the conflict in "conflicting_file"
    When I run `git commit --no-edit; git-town continue`
    Then it runs the commands
      | BRANCH        | COMMAND                      |
      | feature       | git merge --no-edit main     |
      |               | git checkout main            |
      | main          | git merge --squash feature   |
      |               | git commit -m "feature done" |
      |               | git push                     |
      |               | git push origin :feature     |
      |               | git branch -D feature        |
      |               | git checkout other-feature   |
      | other-feature | git stash pop                |
    And I end up on the "other-feature" branch
    And my workspace still contains my uncommitted file
    And there is no "feature" branch
    And my repository still has the following commits
      | BRANCH | LOCATION         | MESSAGE      | FILE NAME        |
      | main   | local and remote | feature done | conflicting_file |
