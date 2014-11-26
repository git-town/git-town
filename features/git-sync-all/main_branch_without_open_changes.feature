Feature: git-sync-all from the main branch

  Background:
    Given I have feature branches named "feature" and "local_feature"
    And my coworker has a feature branch named "remote_feature"
    And the following commits exist in my repository
      | branch         | location         | message               |  file name          |
      | main           | local            | local main commit     | local_main_file     |
      | main           | remote           | remote main commit    | remote_main_file    |
      | feature        | local and remote | feature commit        | feature_file        |
      | local_feature  | local            | local feature commit  | local_feature_file  |
      | remote_feature | remote           | remote feature commit | remote_feature_file |
    And I am on the "main" branch
    When I run `git sync-all`


  Scenario: result
    Then I am still on the "main" branch
    And all branches are now synchronized
    And I have the following commits
      | branch         | location         | message                                | files               |
      | main           | local and remote | local main commit                      | local_main_file     |
      | main           | local and remote | remote main commit                     | remote_main_file    |
      | feature        | local and remote | Merge branch 'main' into feature       |                     |
      | feature        | local and remote | local main commit                      | local_main_file     |
      | feature        | local and remote | remote main commit                     | remote_main_file    |
      | feature        | local and remote | feature commit                         | feature_file        |
      | local_feature  | local and remote | Merge branch 'main' into local_feature |                     |
      | local_feature  | local and remote | local main commit                      | local_main_file     |
      | local_feature  | local and remote | remote main commit                     | remote_main_file    |
      | local_feature  | local and remote | local feature commit                   | local_feature_file  |
      | remote_feature | remote           | remote feature commit                  | remote_feature_file |
