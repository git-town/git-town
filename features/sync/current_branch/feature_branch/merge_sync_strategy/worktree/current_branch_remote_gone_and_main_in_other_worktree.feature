Feature: sync a branch whose remote branch is gone while main is active in another worktree

  Background:
    Given a Git repo with origin
    And the branches
      | NAME       | TYPE    | PARENT    | LOCATIONS     |
      | feature-1  | feature | main      | local, origin |
      | feature-1A | feature | feature-1 | local, origin |
      | feature-2  | feature | main      | local, origin |
    And origin deletes the "feature-1" branch
    And the current branch is "feature-1"
    And branch "main" is active in another worktree
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH     | COMMAND                  |
      | feature-1  | git fetch --prune --tags |
      |            | git checkout feature-1A  |
      | feature-1A | git branch -D feature-1  |
      |            | git checkout feature-2   |
    And this lineage exists now
      """
      main
        feature-1A
        feature-2
      """
    And the initial commits exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                         |
      | feature-2 | git branch feature-1 {{ sha 'initial commit' }} |
      |           | git checkout feature-1                          |
    And this lineage exists now
      """
      main
        feature-1
          feature-1A
        feature-2
      """
    And the initial commits exist now
