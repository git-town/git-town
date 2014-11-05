Feature: git-sync
  on a feature branch
  with two collaborators


  Scenario: merging work
    Given I am on a feature branch
    And my coworker Charlie works on the same feature branch
    And the following commits exist in my repository
      | location  | message     | file name |
      | local     | my commit 1 | my_file_1 |
    And the following commits exist in Charlie's repository
      | location | message           | file name      |
      | local    | charlies commit 1 | charlie_file_1 |
    When I run `git sync`
    Then I see the following commits
      | branch  | location         | message     | files     |
      | feature | local and remote | my commit 1 | my_file_1 |
    And Charlie still sees the following commits
      | branch  | location | message           | files          |
      | feature | local    | charlies commit 1 | charlie_file_1 |
    When Charlie runs `git sync`
    Then now Charlie sees the following commits
      | branch  | location         | message           | files          |
      | feature | local and remote | charlies commit 1 | charlie_file_1 |
      | feature | local and remote | my commit 1       | my_file_1      |
    When I run `git sync`
    Then now I see the following commits
      | branch  | location         | message           | files          |
      | feature | local and remote | my commit 1       | my_file_1      |
      | feature | local and remote | charlies commit 1 | charlie_file_1 |
