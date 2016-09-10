Feature: git-sync: on a feature branch with merge pull branch strategy

  As a developer syncing a feature branch with a merge pull branch strategy
  I want my branch to merge in the main branch instead of rebasing
  So that I never have to rewrite history.


  Background:
    Given my repository has the "pull-branch-strategy" configuration set to "merge"
    And I have a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE               | FILE NAME           |
      | main    | local    | local main commit     | local_main_file     |
      |         | remote   | remote main commit    | remote_main_file    |
      | feature | local    | local feature commit  | local_feature_file  |
      |         | remote   | remote feature commit | remote_feature_file |
    And I am on the "feature" branch
    And I have an uncommitted file
    When I run `git sync`


  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune                  |
      |         | git stash -a                       |
      |         | git checkout main                  |
      | main    | git merge --no-edit origin/main    |
      |         | git push                           |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
      |         | git push                           |
      |         | git stash pop                      |
    And I am still on the "feature" branch
    And I still have my uncommitted file
    And I have the following commits
      | BRANCH  | LOCATION         | MESSAGE                                                    | FILE NAME           |
      | main    | local and remote | local main commit                                          | local_main_file     |
      |         |                  | remote main commit                                         | remote_main_file    |
      |         |                  | Merge remote-tracking branch 'origin/main' into main       |                     |
      | feature | local and remote | local feature commit                                       | local_feature_file  |
      |         |                  | remote feature commit                                      | remote_feature_file |
      |         |                  | Merge remote-tracking branch 'origin/feature' into feature |                     |
      |         |                  | local main commit                                          | local_main_file     |
      |         |                  | remote main commit                                         | remote_main_file    |
      |         |                  | Merge remote-tracking branch 'origin/main' into main       |                     |
      |         |                  | Merge branch 'main' into feature                           |                     |
