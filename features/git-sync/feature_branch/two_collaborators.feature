Feature: Git Sync: collaborative feature branch syncing

  Background:
    Given I have a feature branch named "feature"
    And my coworker fetches updates
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE     | FILE NAME |
      | feature | local    | my commit 1 | my_file_1 |
    And the following commits exist in my coworker's repository
      | BRANCH  | LOCATION | MESSAGE           | FILE NAME      |
      | feature | local    | charlies commit 1 | charlie_file_1 |

  Scenario: result
    And I am on the "feature" branch
    When I run `git sync`
    Then it runs the Git commands
      | BRANCH  | COMMAND                            |
      | feature | git checkout main                  |
      | main    | git fetch --prune                  |
      | main    | git rebase origin/main             |
      | main    | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      | feature | git merge --no-edit main           |
      | feature | git push                           |
    And I have the following commits
      | BRANCH  | LOCATION         | MESSAGE     | FILES     |
      | feature | local and remote | my commit 1 | my_file_1 |

    Given my coworker is on the "feature" branch
    When my coworker runs `git sync`
    Then it runs the Git commands
      | BRANCH  | COMMAND                            |
      | feature | git checkout main                  |
      | main    | git fetch --prune                  |
      | main    | git rebase origin/main             |
      | main    | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      | feature | git merge --no-edit main           |
      | feature | git push                           |
    And now my coworker has the following commits
      | BRANCH  | LOCATION         | MESSAGE                                                    | FILES          |
      | feature | local and remote | Merge remote-tracking branch 'origin/feature' into feature |                |
      | feature |                  | charlies commit 1                                          | charlie_file_1 |
      | feature |                  | my commit 1                                                | my_file_1      |

    Given I am on the "feature" branch
    When I run `git sync`
    Then it runs the Git commands
      | BRANCH  | COMMAND                            |
      | feature | git checkout main                  |
      | main    | git fetch --prune                  |
      | main    | git rebase origin/main             |
      | main    | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      | feature | git merge --no-edit main           |
    And now I have the following commits
      | BRANCH  | LOCATION         | MESSAGE                                                    | FILES          |
      | feature | local and remote | Merge remote-tracking branch 'origin/feature' into feature |                |
      | feature |                  | charlies commit 1                                          | charlie_file_1 |
      | feature |                  | my commit 1                                                | my_file_1      |
