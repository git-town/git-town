Feature: git ship: shipping the supplied feature branch in a fork of the current repo

  (see ../../current_branch/on_feature_branch/without_open_changes/with_tracking_branch.feature)


  Background:
    Given I have a feature branch named "other_feature"
    And a fork of my repository has a branch named "feature"
    And the following commits exist in a fork of my repository
      | BRANCH  | LOCATION | MESSAGE        | FILE NAME    | FILE CONTENT    |
      | feature | remote   | feature commit | feature_file | feature content |
    And I am on the "other_feature" branch


  Scenario: with open changes
    Given I have an uncommitted file with name: "feature_file" and content: "conflicting content"
    When I run `git ship fork:feature -m 'feature done'`
    Then it runs the Git commands
      | BRANCH        | COMMAND                            |
      | other_feature | git stash -u                       |
      | other_feature | git checkout main                  |
      | main          | git fetch --prune                  |
      | main          | git rebase origin/main             |
      | main          | git checkout -b fork/feature       |
      | main          | git remote add fork_origin file:// |
      | main          | git checkout -b fork_origin/feature fork_origin/feature |
      | fork_origin/feature       | git merge --no-edit main           |
      | fork_origin/feature       | git checkout main                  |
      | main          | git merge --squash feature         |
      | main          | git commit -m "feature done"       |
      | main          | git push |
      | main          | git remote remove fork_origin      |
      | main          | git checkout other_feature         |
      | other_feature | git stash pop                      |
    And I end up on the "other_feature" branch
    And I still have an uncommitted file with name: "feature_file" and content: "conflicting content"
    And there is no "feature" branch
    And I have the following commits
      | BRANCH | LOCATION         | MESSAGE      | FILE NAME    |
      | main   | local and remote | feature done | feature_file |


  Scenario: without open changes
    When I run `git ship feature -m "feature done"`
    Then it runs the Git commands
      | BRANCH        | COMMAND                            |
      | other_feature | git checkout main                  |
      | main          | git fetch --prune                  |
      | main          | git rebase origin/main             |
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
    And I end up on the "other_feature" branch
    And there is no "feature" branch
    And I have the following commits
      | BRANCH | LOCATION         | MESSAGE      | FILE NAME    |
      | main   | local and remote | feature done | feature_file |
