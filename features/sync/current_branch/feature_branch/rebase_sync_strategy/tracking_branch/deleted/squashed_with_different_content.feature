Feature: keep syncing a branch whose changes differ from the squash-merge at the remote

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | shipped | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        | FILE NAME | FILE CONTENT   |
      | shipped | local, origin | first commit   | file-a    | first version  |
      | shipped | local, origin | second commit  | file-b    | second version |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the current branch is "shipped"
    And an uncommitted file "file-b" with content "local version"
    And I ran "git commit --all --amend --message 'second commit rewritten'"
    And origin ships the "shipped" branch using the "squash-merge" ship-strategy
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
      CONFLICT (add/add): Merge conflict in file-b
      """
    And a rebase is now in progress
