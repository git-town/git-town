Feature: shipped parent branches in a stacked change

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE    | PARENT | LOCATIONS     |
      | feature-1 | feature | main   | local, origin |
    And the commits
      | BRANCH    | LOCATION      | MESSAGE          | FILE NAME      | FILE CONTENT      |
      | feature-1 | local, origin | feature-1 commit | feature-1-file | feature 1 content |
    And the branches
      | NAME      | TYPE    | PARENT    | LOCATIONS     |
      | feature-2 | feature | feature-1 | local, origin |
    And the commits
      | BRANCH    | LOCATION      | MESSAGE          | FILE NAME      | FILE CONTENT      |
      | feature-2 | local, origin | feature-2 commit | feature-2-file | feature 2 content |
    And the branches
      | NAME      | TYPE    | PARENT    | LOCATIONS     |
      | feature-3 | feature | feature-2 | local, origin |
    And the commits
      | BRANCH    | LOCATION      | MESSAGE          | FILE NAME      | FILE CONTENT      |
      | feature-3 | local, origin | feature-3 commit | feature-3-file | feature 3 content |
    And the branches
      | NAME      | TYPE    | PARENT    | LOCATIONS     |
      | feature-4 | feature | feature-3 | local, origin |
    And the commits
      | BRANCH    | LOCATION      | MESSAGE          | FILE NAME      | FILE CONTENT      |
      | feature-4 | local, origin | feature-4 commit | feature-4-file | feature 4 content |
    And Git Town setting "sync-feature-strategy" is "rebase"
    And origin ships the "feature-1" branch using the "squash-merge" ship-strategy
    And origin ships the "feature-2" branch using the "squash-merge" ship-strategy
    And the current branch is "feature-4"
    When I run "git-town sync"

  @this
  Scenario: result
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                         |
      | feature-4 | git fetch --prune --tags                        |
      |           | git checkout main                               |
      | main      | git rebase origin/main --no-update-refs         |
      |           | git checkout feature-2                          |
      | feature-2 | git rebase --onto main feature-1                |
      |           | git checkout feature-3                          |
      | feature-3 | git pull                                        |
      |           | git rebase --onto main feature-2                |
      |           | git push --force-with-lease                     |
      |           | git rebase feature-2 --no-update-refs           |
      |           | git push --force-with-lease --force-if-includes |
      |           | git checkout feature-4                          |
      | feature-4 | git pull                                        |
      |           | git rebase --onto main feature-2                |
      |           | git push --force-with-lease                     |
      |           | git rebase feature-3 --no-update-refs           |
      |           | git push --force-with-lease --force-if-includes |
      |           | git branch -D feature-1                         |
      |           | git branch -D feature-2                         |
    And Git Town prints:
      """
      deleted branch "feature-1"
      """
    # And Git Town prints:
    #   """
    #   deleted branch "feature-2"
    #   """
    And the current branch is still "feature-4"
    And the branches are now
      | REPOSITORY    | BRANCHES                   |
      | local, origin | main, feature-3, feature-4 |
    And these commits exist now
      | BRANCH    | LOCATION      | MESSAGE          |
      | main      | local, origin | feature-1 commit |
      |           |               | feature-1 commit |
      | feature-3 | local, origin | feature-3 commit |
      | feature-4 | local, origin | feature-4 commit |
    And this lineage exists now
      | BRANCH    | PARENT    |
      | feature-3 | main      |
      | feature-4 | feature-3 |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                                      |
      | feature-3 | git reset --hard {{ sha 'feature-3 commit' }}                |
      |           | git push --force-with-lease --force-if-includes              |
      |           | git checkout main                                            |
      | main      | git reset --hard {{ sha 'initial commit' }}                  |
      |           | git branch feature-1 {{ sha-before-run 'feature-1 commit' }} |
      |           | git branch feature-2 {{ sha-before-run 'feature-2 commit' }} |
      |           | git checkout feature-3                                       |
    And the current branch is still "feature-3"
    And the initial branches and lineage exist now
