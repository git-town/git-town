Feature: git sync: on a feature branch (without open changes)

  (see ./no_conflicts_with_open_changes.feature)


  Scenario: without a remote branch
    Given I have a local feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE              | FILE NAME          |
      | main    | local    | local main commit    | local_main_file    |
      |         | remote   | remote main commit   | remote_main_file   |
      | feature | local    | local feature commit | local_feature_file |
    And I am on the "feature" branch
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
    And all branches are now synchronized
    And I have the following commits
      | BRANCH  | LOCATION         | MESSAGE                          | FILE NAME          |
      | main    | local and remote | local main commit                | local_main_file    |
      |         |                  | remote main commit               | remote_main_file   |
      | feature | local and remote | Merge branch 'main' into feature |                    |
      |         |                  | local main commit                | local_main_file    |
      |         |                  | remote main commit               | remote_main_file   |
      |         |                  | local feature commit             | local_feature_file |


  Scenario: with a remote branch
    Given I have a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE               | FILE NAME           |
      | main    | local    | local main commit     | local_main_file     |
      |         | remote   | remote main commit    | remote_main_file    |
      | feature | local    | local feature commit  | local_feature_file  |
      |         | remote   | remote feature commit | remote_feature_file |
    And I am on the "feature" branch
    When I run `git sync`
    Then it runs the Git commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune                  |
      | feature | git checkout main                  |
      | main    | git rebase origin/main             |
      | main    | git push                           |
      | main    | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      | feature | git merge --no-edit main           |
      | feature | git push                           |
    And I am still on the "feature" branch
    And I have the following commits
      | BRANCH  | LOCATION         | MESSAGE                                                    | FILE NAME           |
      | main    | local and remote | local main commit                                          | local_main_file     |
      |         |                  | remote main commit                                         | remote_main_file    |
      | feature | local and remote | Merge branch 'main' into feature                           |                     |
      |         |                  | Merge remote-tracking branch 'origin/feature' into feature |                     |
      |         |                  | local main commit                                          | local_main_file     |
      |         |                  | remote main commit                                         | remote_main_file    |
      |         |                  | local feature commit                                       | local_feature_file  |
      |         |                  | remote feature commit                                      | remote_feature_file |
