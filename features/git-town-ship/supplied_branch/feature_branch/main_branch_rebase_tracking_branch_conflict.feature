Feature: git town-ship: resolving conflicts between the main branch and its tracking branch

  (see ../../../current_branch/on_feature_branch/without_open_changes/main_branch_rebase_tracking_branch_conflict.feature)


  Background:
    Given my repository has the feature branches "feature" and "other-feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT               |
      | main    | local    | conflicting local commit  | conflicting_file | local conflicting content  |
      |         | remote   | conflicting remote commit | conflicting_file | remote conflicting content |
      | feature | local    | feature commit            | feature_file     | feature content            |
    And I am on the "other-feature" branch
    And my workspace has an uncommitted file
    And I run `git-town ship feature -m "feature done"`


  Scenario: result
    Then it runs the commands
      | BRANCH        | COMMAND                |
      | other-feature | git fetch --prune      |
      |               | git add -A             |
      |               | git stash              |
      |               | git checkout main      |
      | main          | git rebase origin/main |
    And it prints the error:
      """
      To abort, run "git-town ship --abort".
      To continue after you have resolved the conflicts, run "git-town ship --continue".
      """
    And my repo has a rebase in progress
    And my uncommitted file is stashed


  Scenario: aborting
    When I run `git-town ship --abort`
    Then it runs the commands
      | BRANCH        | COMMAND                    |
      | main          | git rebase --abort         |
      |               | git checkout other-feature |
      | other-feature | git stash pop              |
    And I am still on the "other-feature" branch
    And my workspace still contains my uncommitted file
    And there is no rebase in progress
    And my repository is left with my original commits


  Scenario: continuing after resolving the conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `git-town ship --continue`
    Then it runs the commands
      | BRANCH        | COMMAND                            |
      | main          | git rebase --continue              |
      |               | git push                           |
      |               | git checkout feature               |
      | feature       | git merge --no-edit origin/feature |
      |               | git merge --no-edit main           |
      |               | git checkout main                  |
      | main          | git merge --squash feature         |
      |               | git commit -m "feature done"       |
      |               | git push                           |
      |               | git push origin :feature           |
      |               | git branch -D feature              |
      |               | git checkout other-feature         |
      | other-feature | git stash pop                      |
    And I end up on the "other-feature" branch
    And my workspace still contains my uncommitted file
    And there is no "feature" branch
    And my repository still has the following commits
      | BRANCH | LOCATION         | MESSAGE                   | FILE NAME        |
      | main   | local and remote | conflicting remote commit | conflicting_file |
      |        |                  | conflicting local commit  | conflicting_file |
      |        |                  | feature done              | feature_file     |


  Scenario: continuing after resolving the conflicts and continuing the rebase
    Given I resolve the conflict in "conflicting_file"
    When I run `git rebase --continue; git-town ship --continue`
    Then it runs the commands
      | BRANCH        | COMMAND                            |
      | main          | git push                           |
      |               | git checkout feature               |
      | feature       | git merge --no-edit origin/feature |
      |               | git merge --no-edit main           |
      |               | git checkout main                  |
      | main          | git merge --squash feature         |
      |               | git commit -m "feature done"       |
      |               | git push                           |
      |               | git push origin :feature           |
      |               | git branch -D feature              |
      |               | git checkout other-feature         |
      | other-feature | git stash pop                      |
    And I end up on the "other-feature" branch
    And my workspace still contains my uncommitted file
    And there is no "feature" branch
    And my repository still has the following commits
      | BRANCH | LOCATION         | MESSAGE                   | FILE NAME        |
      | main   | local and remote | conflicting remote commit | conflicting_file |
      |        |                  | conflicting local commit  | conflicting_file |
      |        |                  | feature done              | feature_file     |
