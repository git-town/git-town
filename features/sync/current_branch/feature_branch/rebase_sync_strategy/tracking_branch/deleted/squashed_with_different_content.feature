Feature: keep syncing a branch whose changes differ from the squash-merge at the remote

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | shipped | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE   | FILE NAME | FILE CONTENT |
      | shipped | local, origin | commit-1  | file-1    | content 1    |
      | shipped | local, origin | commit-2A | file-2    | content 2    |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And origin ships the "shipped" branch using the "squash-merge" ship-strategy
    And I amend this commit
      | BRANCH  | LOCATION | MESSAGE   | FILE NAME | FILE CONTENT |
      | shipped | local    | commit-2B | file-2    | local change |
    And the current branch is "shipped"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                                      |
      | shipped | git fetch --prune --tags                                                     |
      |         | git checkout main                                                            |
      | main    | git -c rebase.updateRefs=false rebase origin/main                            |
      |         | git checkout shipped                                                         |
      | shipped | git -c rebase.updateRefs=false rebase --onto main {{ sha 'initial commit' }} |
    And Git Town prints the error:
      """
      CONFLICT (add/add): Merge conflict in file-2
      """
    And a rebase is now in progress

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                     |
      | shipped | git rebase --abort                          |
      |         | git checkout main                           |
      | main    | git reset --hard {{ sha 'initial commit' }} |
      |         | git checkout shipped                        |
    And the initial branches and lineage exist now
    And the initial commits exist now
