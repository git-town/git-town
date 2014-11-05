Feature: git-sync
  on a feature branch
  no conflicts


  Scenario: without a remote branch
    Given I am on a local feature branch
    And the following commits exist in my repository
      | branch  | location | message              | file name          |
      | main    | local    | local main commit    | local_main_file    |
      | main    | remote   | remote main commit   | remote_main_file   |
      | feature | local    | local feature commit | local_feature_file |
    When I run `git sync`
    Then I am still on the "feature" branch
    And all branches are now synchronized
    And I have the following commits
      | branch  | location         | message                          | files               |
      | main    | local and remote | local main commit                | local_main_file     |
      | main    | local and remote | remote main commit               | remote_main_file    |
      | feature | local and remote | Merge branch 'main' into feature |                     |
      | feature | local and remote | local main commit                | local_main_file     |
      | feature | local and remote | remote main commit               | remote_main_file    |
      | feature | local and remote | local feature commit             | local_feature_file  |
    And now I have the following committed files
      | branch  | files                                                 |
      | main    | local_main_file, remote_main_file                     |
      | feature | local_feature_file, local_main_file, remote_main_file |


  Scenario: with a remote branch
    Given I am on a feature branch
    And the following commits exist in my repository
      | branch  | location | message               | file name           |
      | main    | local    | local main commit     | local_main_file     |
      | main    | remote   | remote main commit    | remote_main_file    |
      | feature | local    | local feature commit  | local_feature_file  |
      | feature | remote   | remote feature commit | remote_feature_file |
    When I run `git sync`
    Then I am still on the "feature" branch
    And I have the following commits
      | branch  | location         | message                          | files               |
      | main    | local and remote | local main commit                | local_main_file     |
      | main    | local and remote | remote main commit               | remote_main_file    |
      | feature | local and remote | Merge branch 'main' into feature |                     |
      | feature | local and remote | local main commit                | local_main_file     |
      | feature | local and remote | remote main commit               | remote_main_file    |
      | feature | local and remote | local feature commit             | local_feature_file  |
      | feature | local and remote | remote feature commit            | remote_feature_file |
    And now I have the following committed files
      | branch  | files               |
      | main    | local_main_file, remote_main_file |
      | feature | local_feature_file, remote_feature_file, local_main_file, remote_main_file |
