Feature: disable pushing through the CLI

  Scenario: sync all branches with --no-push does not push tags
    Given a Git repo with origin
    And the commits
      | BRANCH | LOCATION      | MESSAGE | FILE NAME | FILE CONTENT |
      | main   | local, origin | commit  | file      | content      |
    And the current branch is "main"
    When I run "git-town sync --all --no-push"
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
    And the initial branches and lineage exist now
    And the initial commits exist now
