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
    And origin ships the "feature-1" branch
    And the current branch is "feature-1" and the previous branch is "feature-2"
    And an uncommitted file
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                 |
      | feature-1 | git fetch --prune --tags                |
      |           | git add -A                              |
      |           | git stash                               |
      |           | git checkout main                       |
      | main      | git rebase origin/main --no-update-refs |
      |           | git branch -D feature-1                 |
      |           | git checkout feature-2                  |
      | feature-2 | git stash pop                           |
    And Git Town prints:
      """
      deleted branch "feature-1"
      """
    And the current branch is now "feature-2"
    And the previous Git branch is now "main"
    And the uncommitted file still exists
    And the branches are now
      | REPOSITORY    | BRANCHES        |
      | local, origin | main, feature-2 |
    And this lineage exists now
      | BRANCH    | PARENT |
      | feature-2 | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                           |
      | feature-2 | git add -A                                        |
      |           | git stash                                         |
      |           | git checkout main                                 |
      | main      | git reset --hard {{ sha 'initial commit' }}       |
      |           | git branch feature-1 {{ sha 'feature-1 commit' }} |
      |           | git checkout feature-1                            |
      | feature-1 | git stash pop                                     |
    And the current branch is now "feature-1"
    And the uncommitted file still exists
    And the initial branches and lineage exist now
