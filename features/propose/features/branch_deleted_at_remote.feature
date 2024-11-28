@skipWindows
Feature: proposing a branch that was deleted toRefId the remote

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And tool "open" is installed
    And the current branch is "feature"
    And the origin is "git@github.com:git-town/git-town.git"
    And origin deletes the "feature" branch
    When I run "git-town propose"

  Scenario: a PR for this branch exists already
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git fetch --prune --tags                |
      |         | git checkout main                       |
      | main    | git rebase origin/main --no-update-refs |
      |         | git branch -D feature                   |
    And Git Town prints:
      """
      branch "feature" was deleted toRefId the remote
      """
    And the current branch is now "main"
    And these branches exist now
      | REPOSITORY    | BRANCHES |
      | local, origin | main     |
    And no lineage exists now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                       |
      | main   | git branch feature {{ sha 'initial commit' }} |
      |        | git checkout feature                          |
    And the current branch is still "feature"
    And the initial branches and lineage exist now
    And the initial commits exist now
