Feature: sync a branch whose tracking branch was shipped in offline mode

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE    | PARENT | LOCATIONS     |
      | feature   | feature | main   | local, origin |
      | feature-1 | feature | main   | local, origin |
      | feature-2 | feature | main   | local, origin |
    And the current branch is "feature"
    And the commits
      | BRANCH    | LOCATION      | MESSAGE          | FILE NAME      | FILE CONTENT      |
      | feature-1 | local, origin | feature-1 commit | feature-1-file | feature 1 content |
      | feature-2 | local, origin | feature-2 commit | feature-2-file | feature 2 content |
    And origin ships the "feature-1" branch
    And the current branch is "feature-1"
    And an uncommitted file
    And offline mode is enabled
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH    | COMMAND                                   |
      | feature-1 | git add -A                                |
      |           | git stash                                 |
      |           | git checkout main                         |
      | main      | git rebase origin/main --no-update-refs   |
      |           | git checkout feature-1                    |
      | feature-1 | git merge --no-edit --ff main             |
      |           | git merge --no-edit --ff origin/feature-1 |
      |           | git stash pop                             |
    And the current branch is still "feature-1"
    And the uncommitted file still exists
    And the initial branches and lineage exist now

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH    | COMMAND       |
      | feature-1 | git add -A    |
      |           | git stash     |
      |           | git stash pop |
    And the current branch is now "feature-1"
    And the uncommitted file still exists
    And the initial branches and lineage exist now
