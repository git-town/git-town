Feature: git sync: syncing the current feature branch (without a tracking branch or remote repo)

  (see ./with_a_tracking_branch.feature)


  Background:
    Given I have a local feature branch named "feature"
    And my repo does not have a remote origin
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE              | FILE NAME          |
      | main    | local    | local main commit    | local_main_file    |
      | feature | local    | local feature commit | local_feature_file |
    And I am on the "feature" branch


  Scenario: with open changes
    Given I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git sync`
    Then it runs the Git commands
      | BRANCH  | COMMAND                  |
      | feature | git stash -u             |
      | feature | git merge --no-edit main |
      | feature | git stash pop            |
    And I am still on the "feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And I have the following commits
      | BRANCH  | LOCATION | MESSAGE                          | FILE NAME          |
      | main    | local    | local main commit                | local_main_file    |
      | feature | local    | local feature commit             | local_feature_file |
      |         |          | local main commit                | local_main_file    |
      |         |          | Merge branch 'main' into feature |                    |


  Scenario: without open changes
    When I run `git sync`
    Then it runs the Git commands
      | BRANCH  | COMMAND                  |
      | feature | git merge --no-edit main |
    And I am still on the "feature" branch
    And I have the following commits
      | BRANCH  | LOCATION | MESSAGE                          | FILE NAME          |
      | main    | local    | local main commit                | local_main_file    |
      | feature | local    | local feature commit             | local_feature_file |
      |         |          | local main commit                | local_main_file    |
      |         |          | Merge branch 'main' into feature |                    |
