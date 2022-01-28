Feature: syncing a nested feature branch (with known parent branches)

  Scenario:
    Given my repo has a feature branch named "parent-feature"
    And my repo has a feature branch named "child-feature" as a child of "parent-feature"
    And the following commits exist in my repo
      | BRANCH         | LOCATION | MESSAGE                      | FILE NAME                  |
      | main           | local    | local main commit            | local_main_file            |
      |                | remote   | remote main commit           | remote_main_file           |
      | parent-feature | local    | local parent feature commit  | local_parent_feature_file  |
      |                | remote   | remote parent feature commit | remote_parent_feature_file |
      | child-feature  | local    | local child feature commit   | local_child_feature_file   |
      |                | remote   | remote child feature commit  | remote_child_feature_file  |
    And I am on the "child-feature" branch
    And my workspace has an uncommitted file
    When I run "git-town sync"
    Then it runs the commands
      | BRANCH         | COMMAND                                   |
      | child-feature  | git fetch --prune --tags                  |
      |                | git add -A                                |
      |                | git stash                                 |
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
    And my workspace still contains my uncommitted file
    And my repo now has the following commits
      | BRANCH         | LOCATION      | MESSAGE                                                                  | FILE NAME                  |
      | main           | local, remote | remote main commit                                                       | remote_main_file           |
      |                |               | local main commit                                                        | local_main_file            |
      | child-feature  | local, remote | local child feature commit                                               | local_child_feature_file   |
      |                |               | remote child feature commit                                              | remote_child_feature_file  |
      |                |               | Merge remote-tracking branch 'origin/child-feature' into child-feature   |                            |
      |                |               | local parent feature commit                                              | local_parent_feature_file  |
      |                |               | remote parent feature commit                                             | remote_parent_feature_file |
      |                |               | Merge remote-tracking branch 'origin/parent-feature' into parent-feature |                            |
      |                |               | remote main commit                                                       | remote_main_file           |
      |                |               | local main commit                                                        | local_main_file            |
      |                |               | Merge branch 'main' into parent-feature                                  |                            |
      |                |               | Merge branch 'parent-feature' into child-feature                         |                            |
      | parent-feature | local, remote | local parent feature commit                                              | local_parent_feature_file  |
      |                |               | remote parent feature commit                                             | remote_parent_feature_file |
      |                |               | Merge remote-tracking branch 'origin/parent-feature' into parent-feature |                            |
      |                |               | remote main commit                                                       | remote_main_file           |
      |                |               | local main commit                                                        | local_main_file            |
      |                |               | Merge branch 'main' into parent-feature                                  |                            |
