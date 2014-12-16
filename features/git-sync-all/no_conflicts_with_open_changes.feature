Feature: git-sync-all from the main branch

  Background:
    Given I have feature branches named "feature1", "feature2", and "feature3"
    And my coworker has a feature branch named "coworker_feature"
    And the following commits exist in my repository
      | branch           | location         | message                | file name            |
      | main             | local            | main local commit      | main_local_file      |
      | main             | remote           | main remote commit     | main_remote_file     |
      | feature1         | local and remote | feature1 commit        | feature1_file        |
      | feature2         | local            | feature2 local commit  | feature2_local_file  |
      |                  | remote           | feature2 remote commit | feature2_remote_file |
      | feature3         | local            | feature3 commit        | feature3_file        |
      | coworker_feature | remote           | coworker commit        | coworker_file        |
    And I am on the "main" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git sync-all`


  Scenario: result
    Then I am still on the "main" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And all branches are now synchronized
    And I have the following commits
      | branch           | location         | message                                                      | files                |
      | main             | local and remote | main local commit                                            | main_local_file      |
      | main             | local and remote | main remote commit                                           | main_remote_file     |
      | feature1         | local and remote | Merge branch 'main' into feature1                            |                      |
      |                  |                  | main local commit                                            | main_local_file      |
      |                  |                  | main remote commit                                           | main_remote_file     |
      |                  |                  | feature1 commit                                              | feature1_file        |
      | feature2         | local and remote | Merge branch 'main' into feature2                            |                      |
      |                  |                  | Merge remote-tracking branch 'origin/feature2' into feature2 |                      |
      |                  |                  | main local commit                                            | main_local_file      |
      |                  |                  | main remote commit                                           | main_remote_file     |
      |                  |                  | feature2 local commit                                        | feature2_local_file  |
      |                  |                  | feature2 remote commit                                       | feature2_remote_file |
      | feature3         | local and remote | Merge branch 'main' into feature3                            |                      |
      |                  |                  | main local commit                                            | main_local_file      |
      |                  |                  | main remote commit                                           | main_remote_file     |
      |                  |                  | feature3 commit                                              | feature3_file        |
      | coworker_feature | remote           | coworker commit                                              | coworker_file        |
