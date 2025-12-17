Feature: sync a branch whose tracking branch was shipped

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE    | PARENT | LOCATIONS     |
      | feature-1 | feature | main   | local, origin |
      | feature-2 | feature | main   | local, origin |
      | feature-3 | feature | main   | local, origin |
    And the commits
      | BRANCH    | LOCATION      | MESSAGE          |
      | feature-1 | local, origin | feature-1 commit |
      | feature-2 | local, origin | feature-2 commit |
      | feature-3 | local, origin | feature-3 commit |
    And origin ships the "feature-1" branch using the "squash-merge" ship-strategy
    And the current branch is "feature-2"
    And branch "main" is active in another worktree
    When I run "git-town sync --all"

  @this
  Scenario: result
    Then Git Town runs the commands
      | BRANCH    | COMMAND                  |
      | feature-2 | git fetch --prune --tags |
      |           | git branch -D feature-1  |
      |           | git checkout feature-3   |
      | feature-3 | git checkout feature-2   |
      | feature-2 | git push --tags          |
    And Git Town prints:
      """
      deleted branch "feature-1"
      """
    And this lineage exists now
      """
      main
        feature-2
        feature-3
      """
    And the branches are now
      | REPOSITORY    | BRANCHES                   |
      | local, origin | main, feature-2, feature-3 |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                         |
      | main   | git reset --hard {{ sha 'initial commit' }}     |
      |        | git branch feature {{ sha 'feature-1 commit' }} |
      |        | git checkout feature                            |
    And the initial branches and lineage exist now
