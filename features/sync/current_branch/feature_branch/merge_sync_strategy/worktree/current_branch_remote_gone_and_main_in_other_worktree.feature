Feature: sync a branch whose branch is gone while main is active in another worktree

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE    | PARENT | LOCATIONS     |
      | feature-1 | feature | main   | local, origin |
      | feature-2 | feature | main   | local, origin |
    And the current branch is "feature-1"
    And origin deletes the "feature-1" branch
    And branch "main" is active in another worktree
    When I run "git-town sync"

  @this
  Scenario: result
    Then Git Town runs the commands
      | BRANCH    | COMMAND                  |
      | feature-1 | git fetch --prune --tags |
      |           | git checkout feature-2   |
      | feature-2 | git branch -D feature-1  |
    And the initial commits exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                                                |
      | feature | git reset --hard {{ sha 'local feature commit' }}                                      |
      |         | git push --force-with-lease origin {{ sha-in-origin 'origin feature commit' }}:feature |
    And the initial commits exist now
