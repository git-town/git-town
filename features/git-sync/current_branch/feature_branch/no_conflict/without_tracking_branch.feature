Feature: git sync: syncing the current feature branch without a tracking branch

  (see ./with_a_tracking_branch.feature)


  Background:
    Given I have a local feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE              | FILE NAME          |
      | main    | local    | local main commit    | local_main_file    |
      |         | remote   | remote main commit   | remote_main_file   |
      | feature | local    | local feature commit | local_feature_file |
    And I am on the "feature" branch


  Scenario: with open changes
    Given I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git sync`
    Then it runs the Git commands
      | BRANCH  | COMMAND                    |
      | feature | git fetch --prune          |
      | feature | git stash -u               |
      | feature | git checkout main          |
      | main    | git rebase origin/main     |
      | main    | git push                   |
      | main    | git checkout feature       |
      | feature | git merge --no-edit main   |
      | feature | git push -u origin feature |
      | feature | git stash pop              |
    And I am still on the "feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And I have the following commits
      | BRANCH  | LOCATION         | MESSAGE                          | FILE NAME          |
      | main    | local and remote | remote main commit               | remote_main_file   |
      |         |                  | local main commit                | local_main_file    |
      | feature | local and remote | local feature commit             | local_feature_file |
      |         |                  | remote main commit               | remote_main_file   |
      |         |                  | local main commit                | local_main_file    |
      |         |                  | Merge branch 'main' into feature |                    |


  Scenario: without open changes
    When I run `git sync`
    Then it runs the Git commands
      | BRANCH  | COMMAND                    |
      | feature | git fetch --prune          |
      | feature | git checkout main          |
      | main    | git rebase origin/main     |
      | main    | git push                   |
      | main    | git checkout feature       |
      | feature | git merge --no-edit main   |
      | feature | git push -u origin feature |
    And I am still on the "feature" branch
    And I have the following commits
      | BRANCH  | LOCATION         | MESSAGE                          | FILE NAME          |
      | main    | local and remote | remote main commit               | remote_main_file   |
      |         |                  | local main commit                | local_main_file    |
      | feature | local and remote | local feature commit             | local_feature_file |
      |         |                  | remote main commit               | remote_main_file   |
      |         |                  | local main commit                | local_main_file    |
      |         |                  | Merge branch 'main' into feature |                    |
