Feature: sync a branch shipped at the remote, fully shipped, local commits have different metadata

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | shipped | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE   | FILE NAME | FILE CONTENT |
      | shipped | local, origin | commit-1  | file      | content 1    |
      | shipped | local, origin | commit-2A | file      | content 2    |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And origin ships the "shipped" branch using the "squash-merge" ship-strategy
    And the current branch is "shipped"
    And I ran "git commit --amend --message 'commit-2B'"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                           |
      | shipped | git fetch --prune --tags                          |
      |         | git checkout main                                 |
      | main    | git -c rebase.updateRefs=false rebase origin/main |
      |         | git checkout shipped                              |
      | shipped | git checkout main                                 |
      | main    | git branch -D shipped                             |
    And no rebase is now in progress
    And the branches are now
      | REPOSITORY    | BRANCHES |
      | local, origin | main     |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE  |
      | main   | local, origin | commit-1 |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                     |
      | main   | git reset --hard {{ sha 'initial commit' }} |
      |        | git branch shipped {{ sha 'commit-2B' }}    |
      |        | git checkout shipped                        |
    And the initial branches and lineage exist now
    And the initial commits exist now
