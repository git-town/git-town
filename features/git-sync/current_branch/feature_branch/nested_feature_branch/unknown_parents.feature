Feature: git sync: syncing a nested feature branch (without known parent branches)

  As a developer developing a feature branch that was cut from another feature branch
  I want my branch to be synced off the updated parent branch
  So that I get correct updates when I sync my feature branch.


  Scenario: entering the branch names
    Given I have a feature branch named "parent-feature"
    And I have a feature branch named "child-feature" that is cut from parent-feature
    And Git Town has no branch hierarchy information for "parent-feature" and "child-feature"
    And the following commits exist in my repository
      | BRANCH         | LOCATION | MESSAGE                      | FILE NAME                  |
      | main           | local    | local main commit            | local_main_file            |
      |                | remote   | remote main commit           | remote_main_file           |
      | parent-feature | local    | local parent feature commit  | local_parent_feature_file  |
      |                | remote   | remote parent feature commit | remote_parent_feature_file |
      | child-feature  | local    | local child feature commit   | local_child_feature_file   |
      |                | remote   | remote child feature commit  | remote_child_feature_file  |
    And I am on the "child-feature" branch
    And I have an uncommitted file
    When I run `git sync` and enter "parent-feature" and "main"
    Then it runs the Git commands
      | BRANCH         | COMMAND                                   |
      | child-feature  | git fetch --prune                         |
      |                | git stash -u                              |
      |                | git checkout main                         |
      | main           | git rebase origin/main                    |
      |                | git push                                  |
      |                | git checkout parent-feature               |
      | parent-feature | git merge --no-edit origin/parent-feature |
      |                | git merge --no-edit main                  |
      |                | git push                                  |
      |                | git checkout child-feature                |
      | child-feature  | git merge --no-edit origin/child-feature  |
      |                | git merge --no-edit parent-feature        |
      |                | git push                                  |
      |                | git stash pop                             |
    And I am still on the "child-feature" branch
    And I still have my uncommitted file
    And I have the following commits
      | BRANCH         | LOCATION         | MESSAGE                                                                  | FILE NAME                  |
      | main           | local and remote | remote main commit                                                       | remote_main_file           |
      |                |                  | local main commit                                                        | local_main_file            |
      | child-feature  | local and remote | local child feature commit                                               | local_child_feature_file   |
      |                |                  | remote child feature commit                                              | remote_child_feature_file  |
      |                |                  | Merge remote-tracking branch 'origin/child-feature' into child-feature   |                            |
      |                |                  | local parent feature commit                                              | local_parent_feature_file  |
      |                |                  | remote parent feature commit                                             | remote_parent_feature_file |
      |                |                  | Merge remote-tracking branch 'origin/parent-feature' into parent-feature |                            |
      |                |                  | remote main commit                                                       | remote_main_file           |
      |                |                  | local main commit                                                        | local_main_file            |
      |                |                  | Merge branch 'main' into parent-feature                                  |                            |
      |                |                  | Merge branch 'parent-feature' into child-feature                         |                            |
      | parent-feature | local and remote | local parent feature commit                                              | local_parent_feature_file  |
      |                |                  | remote parent feature commit                                             | remote_parent_feature_file |
      |                |                  | Merge remote-tracking branch 'origin/parent-feature' into parent-feature |                            |
      |                |                  | remote main commit                                                       | remote_main_file           |
      |                |                  | local main commit                                                        | local_main_file            |
      |                |                  | Merge branch 'main' into parent-feature                                  |                            |


  Scenario: choosing defaults for the branch names
    Given I have a feature branch named "feature"
    And Git Town has no branch hierarchy information for "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE               | FILE NAME           |
      | main    | local    | local main commit     | local_main_file     |
      |         | remote   | remote main commit    | remote_main_file    |
      | feature | local    | local feature commit  | local_feature_file  |
      |         | remote   | remote feature commit | remote_feature_file |
    And I am on the "feature" branch
    And I have an uncommitted file
    When I run `git sync` and enter ""
    Then it runs the Git commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune                  |
      |         | git stash -u                       |
      |         | git checkout main                  |
      | main    | git rebase origin/main             |
      |         | git push                           |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
      |         | git push                           |
      |         | git stash pop                      |
    And I am still on the "feature" branch
    And I still have my uncommitted file
    And I have the following commits
      | BRANCH  | LOCATION         | MESSAGE                                                    | FILE NAME                  |
      | main    | local and remote | remote main commit                                         | remote_main_file           |
      |         |                  | local main commit                                          | local_main_file            |
      | feature | local and remote | local feature commit                                       | local_feature_file  |
      |         |                  | remote feature commit                                      | remote_feature_file |
      |         |                  | Merge remote-tracking branch 'origin/feature' into feature |                            |
      |         |                  | remote main commit                                         | remote_main_file           |
      |         |                  | local main commit                                          | local_main_file            |
      |         |                  | Merge branch 'main' into feature                           |                            |

