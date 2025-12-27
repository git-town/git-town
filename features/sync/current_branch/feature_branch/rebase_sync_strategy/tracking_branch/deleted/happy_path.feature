@smoke
Feature: sync a branch whose tracking branch was shipped

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE    | PARENT | LOCATIONS     |
      | feature-1 | feature | main   | local, origin |
      | feature-2 | feature | main   | local, origin |
    And the commits
      | BRANCH    | LOCATION      | MESSAGE          | FILE NAME      | FILE CONTENT      |
      | feature-1 | local, origin | feature-1 commit | feature-1-file | feature 1 content |
      | feature-2 | local, origin | feature-2 commit | feature-2-file | feature 2 content |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And origin ships the "feature-1" branch using the "squash-merge" ship-strategy
    And the current branch is "feature-1"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                           |
      | feature-1 | git fetch --prune --tags                          |
      |           | git checkout main                                 |
      | main      | git -c rebase.updateRefs=false rebase origin/main |
      |           | git branch -D feature-1                           |
      |           | git checkout feature-2                            |
    And Git Town prints:
      """
      deleted branch feature-1
      """
    And this lineage exists now
      """
      main
        feature-2
      """
    And the branches are now
      | REPOSITORY    | BRANCHES        |
      | local, origin | main, feature-2 |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                           |
      | feature-2 | git checkout main                                 |
      | main      | git reset --hard {{ sha 'initial commit' }}       |
      |           | git branch feature-1 {{ sha 'feature-1 commit' }} |
      |           | git checkout feature-1                            |
    And the initial branches and lineage exist now
