@skipWindows
Feature: proposing a branch that was deleted at the remote

  Background: proposing changes
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And tool "open" is installed
    And the current branch is "feature"
    And the origin is "git@github.com:git-town/git-town.git"
    And origin deletes the "feature" branch
    When I run "git-town propose"

  @this
  Scenario: a PR for this branch exists already
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                            |
      | feature | git fetch --prune --tags                                           |
      |         | git checkout main                                                  |
      | main    | git rebase origin/main --no-update-refs                            |
      |         | git branch -D feature                                              |
      | <none>  | open https://github.com/git-town/git-town/compare/feature?expand=1 |
    And Git Town prints the error:
      """
      branch "feature" was deleted at the remote
      """
    And the current branch is now "main"
    And the initial branches and lineage exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the current branch is still "feature"
    And the initial commits exist now
