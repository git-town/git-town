Feature: Git Sync: collaborative feature branch syncing



  Background:
    Given I am on a feature branch
    And my coworker Charlie works on the same feature branch
    And the following commits exist in my repository
      | LOCATION  | MESSAGE     | FILE NAME |
      | local     | my commit 1 | my_file_1 |
    And the following commits exist in Charlie's repository
      | LOCATION | MESSAGE           | FILE NAME      |
      | local    | charlies commit 1 | charlie_file_1 |
    When I run `git sync`


  Scenario: result
    Then I have the following commits
      | BRANCH  | LOCATION         | MESSAGE     | FILES     |
      | feature | local and remote | my commit 1 | my_file_1 |
    And Charlie still has the following commits
      | BRANCH  | LOCATION | MESSAGE           | FILES          |
      | feature | local    | charlies commit 1 | charlie_file_1 |
    When Charlie runs `git sync`
    Then now Charlie has the following commits
      | BRANCH  | LOCATION         | MESSAGE                                                    | FILES          |
      | feature | local and remote | Merge remote-tracking branch 'origin/feature' into feature |                |
      | feature |                  | charlies commit 1                                          | charlie_file_1 |
      | feature |                  | my commit 1                                                | my_file_1      |
    When I run `git sync`
    Then now I have the following commits
      | BRANCH  | LOCATION         | MESSAGE                                                    | FILES          |
      | feature | local and remote | Merge remote-tracking branch 'origin/feature' into feature |                |
      | feature |                  | charlies commit 1                                          | charlie_file_1 |
      | feature |                  | my commit 1                                                | my_file_1      |
