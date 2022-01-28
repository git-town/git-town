Feature: syncing the current feature branch without a tracking branch

  Background:
    Given my repo has a local feature branch named "feature"
    And the following commits exist in my repo
      | BRANCH  | LOCATION | MESSAGE              | FILE NAME          |
      | main    | local    | local main commit    | local_main_file    |
      |         | remote   | remote main commit   | remote_main_file   |
      | feature | local    | local feature commit | local_feature_file |
    And I am on the "feature" branch
    And my workspace has an uncommitted file
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                    |
      | feature | git fetch --prune --tags   |
      |         | git add -A                 |
      |         | git stash                  |
      |         | git checkout main          |
      | main    | git rebase origin/main     |
      |         | git push                   |
      |         | git checkout feature       |
      | feature | git merge --no-edit main   |
      |         | git push -u origin feature |
      |         | git stash pop              |
    And I am still on the "feature" branch
    And my workspace still contains my uncommitted file
    And my repo now has the following commits
      | BRANCH  | LOCATION      | MESSAGE                          | FILE NAME          |
      | main    | local, remote | remote main commit               | remote_main_file   |
      |         |               | local main commit                | local_main_file    |
      | feature | local, remote | local feature commit             | local_feature_file |
      |         |               | remote main commit               | remote_main_file   |
      |         |               | local main commit                | local_main_file    |
      |         |               | Merge branch 'main' into feature |                    |
