Feature: git-sync: on a feature branch with merge pull branch strategy

  Background:
    Given the pull-branch-strategy configuration is "merge"
    And my repo has a feature branch named "feature"
    And the following commits exist in my repo
      | BRANCH  | LOCATION | MESSAGE               | FILE NAME           |
      | main    | local    | local main commit     | local_main_file     |
      |         | remote   | remote main commit    | remote_main_file    |
      | feature | local    | local feature commit  | local_feature_file  |
      |         | remote   | remote feature commit | remote_feature_file |
    And I am on the "feature" branch
    And my workspace has an uncommitted file
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune --tags           |
      |         | git add -A                         |
      |         | git stash                          |
      |         | git checkout main                  |
      | main    | git merge --no-edit origin/main    |
      |         | git push                           |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
      |         | git push                           |
      |         | git stash pop                      |
    And I am still on the "feature" branch
    And my workspace still contains my uncommitted file
    And my repo now has the following commits
      | BRANCH  | LOCATION      | MESSAGE                                                    | FILE NAME           |
      | main    | local, remote | local main commit                                          | local_main_file     |
      |         |               | remote main commit                                         | remote_main_file    |
      |         |               | Merge remote-tracking branch 'origin/main'                 |                     |
      | feature | local, remote | local feature commit                                       | local_feature_file  |
      |         |               | remote feature commit                                      | remote_feature_file |
      |         |               | Merge remote-tracking branch 'origin/feature' into feature |                     |
      |         |               | local main commit                                          | local_main_file     |
      |         |               | remote main commit                                         | remote_main_file    |
      |         |               | Merge remote-tracking branch 'origin/main'                 |                     |
      |         |               | Merge branch 'main' into feature                           |                     |
