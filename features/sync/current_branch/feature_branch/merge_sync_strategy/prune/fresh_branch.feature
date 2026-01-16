Feature: prune a freshly created branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    And the current branch is "feature"
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
