Feature: git-town sync: offline mode

  Background:
    Given Git Town is in offline mode
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
      | feature | git add -A                         |
      |         | git stash                          |
      |         | git checkout main                  |
      | main    | git rebase origin/main             |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
      |         | git stash pop                      |
    And I am still on the "feature" branch
    And my workspace still contains my uncommitted file
    And my repo now has the following commits
      | BRANCH  | LOCATION | MESSAGE                          | FILE NAME           |
      | main    | local    | local main commit                | local_main_file     |
      |         | remote   | remote main commit               | remote_main_file    |
      | feature | local    | local feature commit             | local_feature_file  |
      |         |          | local main commit                | local_main_file     |
      |         |          | Merge branch 'main' into feature |                     |
      |         | remote   | remote feature commit            | remote_feature_file |
