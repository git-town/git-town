Feature: leave independent modify/modify conflicts alone when pruning an empty branch

  Background:
    Given a Git repo with origin
    And the commits
      | BRANCH | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT |
      | main   | local, origin | base commit | file      | base content |
    And the branches
      | NAME      | TYPE    | PARENT    | LOCATIONS     |
      | feature-1 | feature | main      | local, origin |
      | feature-2 | feature | feature-1 | local, origin |
    And the commits
      | BRANCH    | LOCATION      | MESSAGE          | FILE NAME | FILE CONTENT      |
      | feature-2 | local, origin | feature-2 commit | file      | feature-2 content |
      | main      | local, origin | main commit      | file      | main content      |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the current branch is "feature-1"
    When I run "git-town sync --prune"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                                                   |
      | feature-1 | git fetch --prune --tags                                                  |
      |           | git -c rebase.updateRefs=false rebase --onto main {{ sha 'base commit' }} |
      |           | git checkout feature-2                                                    |
      | feature-2 | git push origin :feature-1                                                |
      |           | git -c rebase.updateRefs=false rebase --onto main feature-1               |
    And Git Town prints the error:
      """
      CONFLICT (content): Merge conflict in file
      """
    And a rebase is now in progress
    And the branches are now
      | REPOSITORY | BRANCHES                   |
      | local      | main, feature-1, feature-2 |
      | origin     | main, feature-2            |
