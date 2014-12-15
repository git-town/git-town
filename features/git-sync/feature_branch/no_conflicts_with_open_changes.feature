eature: git sync: on a feature branch (with open changes)

  As a developer working on a feature
  I want to be able to easily update my feature branch to include the latest changes from the rest of the team
  So that my work stays in sync with the main development line, can be merged easily later, and I remain productive.


  Scenario: without a remote branch
    Given I am on a local feature branch
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE              | FILE NAME          |
      | main    | local    | local main commit    | local_main_file    |
      |         | remote   | remote main commit   | remote_main_file   |
      | feature | local    | local feature commit | local_feature_file |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git sync`
    Then I am still on the "feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And all branches are now synchronized
    And I have the following commits
      | BRANCH  | LOCATION         | MESSAGE                          | FILES              |
      | main    | local and remote | local main commit                | local_main_file    |
      |         |                  | remote main commit               | remote_main_file   |
      | feature | local and remote | Merge branch 'main' into feature |                    |
      |         |                  | local main commit                | local_main_file    |
      |         |                  | remote main commit               | remote_main_file   |
      |         |                  | local feature commit             | local_feature_file |
    And now I have the following committed files
      | BRANCH  | FILES                                                 |
      | main    | local_main_file, remote_main_file                     |
      | feature | local_feature_file, local_main_file, remote_main_file |


  Scenario: with a remote branch
    Given I am on a feature branch
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE               | FILE NAME           |
      | main    | local    | local main commit     | local_main_file     |
      |         | remote   | remote main commit    | remote_main_file    |
      | feature | local    | local feature commit  | local_feature_file  |
      |         | remote   | remote feature commit | remote_feature_file |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git sync`
    Then I am still on the "feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And I have the following commits
      | BRANCH  | LOCATION         | MESSAGE                                                    | FILES               |
      | main    | local and remote | local main commit                                          | local_main_file     |
      |         |                  | remote main commit                                         | remote_main_file    |
      | feature | local and remote | Merge branch 'main' into feature                           |                     |
      |         |                  | Merge remote-tracking branch 'origin/feature' into feature |                     |
      |         |                  | local main commit                                          | local_main_file     |
      |         |                  | remote main commit                                         | remote_main_file    |
      |         |                  | local feature commit                                       | local_feature_file  |
      |         |                  | remote feature commit                                      | remote_feature_file |
    And now I have the following committed files
      | BRANCH  | FILES               |
      | main    | local_main_file, remote_main_file |
      | feature | local_feature_file, remote_feature_file, local_main_file, remote_main_file |
