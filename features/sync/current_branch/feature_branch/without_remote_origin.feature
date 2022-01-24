Feature: git-town sync: syncing the current feature branch (without a tracking branch or remote repo)

  (see ./with_a_tracking_branch.feature)

  Background:
    Given my repo does not have a remote origin
    And my repo has a local feature branch named "feature"
    And the following commits exist in my repo
      | BRANCH  | LOCATION | MESSAGE              | FILE NAME          | FILE CONTENT    |
      | main    | local    | local main commit    | local_main_file    | main content    |
      | feature | local    | local feature commit | local_feature_file | feature content |
    And I am on the "feature" branch
    And my workspace has an uncommitted file
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git add -A               |
      |         | git stash                |
      |         | git merge --no-edit main |
      |         | git stash pop            |
    And I am still on the "feature" branch
    And my workspace still contains my uncommitted file
    And my repo now has the following commits
      | BRANCH  | LOCATION | MESSAGE                          | FILE NAME          |
      | main    | local    | local main commit                | local_main_file    |
      | feature | local    | local feature commit             | local_feature_file |
      |         |          | local main commit                | local_main_file    |
      |         |          | Merge branch 'main' into feature |                    |
    And my repo now has the following committed files
      | BRANCH  | NAME               | CONTENT         |
      | main    | local_main_file    | main content    |
      | feature | local_feature_file | feature content |
      | feature | local_main_file    | main content    |
