Feature: sync a branch whose tracking branch was shipped

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE          | FILE NAME      | FILE CONTENT      |
      | feature | local, origin | feature-1 commit | feature-1-file | feature 1 content |
    And origin ships the "feature" branch using the squash-merge ship-strategy
    And the current branch is "feature"
    When I run "git-town sync --all"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git fetch --prune --tags                |
      |         | git checkout main                       |
      | main    | git rebase origin/main --no-update-refs |
      |         | git branch -D feature                   |
      |         | git push --tags                         |
    And Git Town prints:
      """
      deleted branch "feature"
      """
    And the current branch is now "main"
    And the branches are now
      | REPOSITORY    | BRANCHES |
      | local, origin | main     |
    And no lineage exists now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                         |
      | main   | git reset --hard {{ sha 'initial commit' }}     |
      |        | git branch feature {{ sha 'feature-1 commit' }} |
      |        | git checkout feature                            |
    And the current branch is now "feature"
    And the initial branches and lineage exist now
