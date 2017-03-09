Feature: git town-sync: syncing the current feature branch without a tracking branch

  (see ./with_a_tracking_branch.feature)


  Background:
    Given I have a local feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE              | FILE NAME          |
      | main    | local    | local main commit    | local_main_file    |
      |         | remote   | remote main commit   | remote_main_file   |
      | feature | local    | local feature commit | local_feature_file |
    And I am on the "feature" branch
    And I have an uncommitted file
    When I run `gt sync`


  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                    |
      | feature | git fetch --prune          |
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
    And I still have my uncommitted file
    And I have the following commits
      | BRANCH  | LOCATION         | MESSAGE                          | FILE NAME          |
      | main    | local and remote | remote main commit               | remote_main_file   |
      |         |                  | local main commit                | local_main_file    |
      | feature | local and remote | local feature commit             | local_feature_file |
      |         |                  | remote main commit               | remote_main_file   |
      |         |                  | local main commit                | local_main_file    |
      |         |                  | Merge branch 'main' into feature |                    |
