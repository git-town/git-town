Feature: git ship: resolving main branch updates when shipping a given feature branch (with open changes)

  (see ../current_branch/pull_main_branch_conflict.feature)


  Background:
    Given I have feature branches named "feature" and "other_feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT               |
      | main    | remote   | conflicting remote commit | conflicting_file | remote conflicting content |
      |         | local    | conflicting local commit  | conflicting_file | local conflicting content  |
      | feature | local    | feature commit            | feature_file     | feature content            |
    And I am on the "other_feature" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    And I run `git ship feature -m "feature done"` while allowing errors


  @finishes-with-non-empty-stash
  Scenario: result
    Then it runs the Git commands
      | BRANCH        | COMMAND                |
      | other_feature | git stash -u           |
      | other_feature | git checkout main      |
      | main          | git fetch --prune      |
      | main          | git rebase origin/main |
    And my repo has a rebase in progress
    And I don't have an uncommitted file with name: "uncommitted"


  Scenario: aborting
    When I run `git ship --abort`
    Then it runs the Git commands
      | BRANCH        | COMMAND                    |
      | HEAD          | git rebase --abort         |
      | main          | git checkout other_feature |
      | other_feature | git stash pop              |
    And I am still on the "other_feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And there is no rebase in progress
    And I am left with my original commits


  Scenario: continuing after resolving conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `git ship --continue`
    Then it runs the Git commands
      | BRANCH        | COMMAND                            |
      | HEAD          | git rebase --continue              |
      | main          | git push                           |
      | main          | git checkout feature               |
      | feature       | git merge --no-edit origin/feature |
      | feature       | git merge --no-edit main           |
      | feature       | git checkout main                  |
      | main          | git merge --squash feature         |
      | main          | git commit -m "feature done"       |
      | main          | git push                           |
      | main          | git push origin :feature           |
      | main          | git branch -D feature              |
      | main          | git checkout other_feature         |
      | other_feature | git stash pop                      |
    And I end up on the "other_feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And there is no "feature" branch
    And I still have the following commits
      | BRANCH | LOCATION         | MESSAGE                   | FILE NAME        |
      | main   | local and remote | feature done              | feature_file     |
      |        |                  | conflicting local commit  | conflicting_file |
      |        |                  | conflicting remote commit | conflicting_file |


  Scenario: continuing after resolving conflicts and continuing the rebase
    Given I resolve the conflict in "conflicting_file"
    When I run `git rebase --continue; git ship --continue`
    Then it runs the Git commands
      | BRANCH        | COMMAND                            |
      | main          | git push                           |
      | main          | git checkout feature               |
      | feature       | git merge --no-edit origin/feature |
      | feature       | git merge --no-edit main           |
      | feature       | git checkout main                  |
      | main          | git merge --squash feature         |
      | main          | git commit -m "feature done"       |
      | main          | git push                           |
      | main          | git push origin :feature           |
      | main          | git branch -D feature              |
      | main          | git checkout other_feature         |
      | other_feature | git stash pop                      |
    And I end up on the "other_feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And there is no "feature" branch
    And I still have the following commits
      | BRANCH | LOCATION         | MESSAGE                   | FILE NAME        |
      | main   | local and remote | feature done              | feature_file     |
      |        |                  | conflicting local commit  | conflicting_file |
      |        |                  | conflicting remote commit | conflicting_file |
