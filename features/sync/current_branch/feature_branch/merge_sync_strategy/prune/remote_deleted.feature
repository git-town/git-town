Feature: prune a branch with deleted remote

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        | FILE NAME    |
      | feature | local, origin | feature commit | feature_file |
    And the current branch is "feature"
    And origin deletes the "feature" branch
    When I run "git-town sync --prune"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git fetch --prune --tags |
      |         | git checkout main        |
      | main    | git branch -D feature    |
    And the branches are now
      | REPOSITORY    | BRANCHES |
      | local, origin | main     |
    And all branches are now synchronized

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                       |
      | main   | git branch feature {{ sha 'initial commit' }} |
      |        | git checkout feature                          |
    And the initial branches and lineage exist now
    And the initial commits exist now
