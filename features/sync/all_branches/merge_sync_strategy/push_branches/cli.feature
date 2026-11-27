Feature: disable pushing through the CLI

  Background:
    Given a Git repo with origin
    And the commits
      | BRANCH | LOCATION | MESSAGE     |
      | main   | origin   | main commit |
    And the current branch is "main"
    And the tags
      | NAME      | LOCATION |
      | local-tag | local    |
    When I run "git-town sync --all --no-push"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                           |
      | main   | git fetch --prune --tags                          |
      |        | git -c rebase.updateRefs=false rebase origin/main |
    And the initial branches and lineage exist now
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
    And the initial tags exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                     |
      | main   | git reset --hard {{ sha 'initial commit' }} |
    And the initial branches and lineage exist now
    And the initial commits exist now
    And the initial tags exist now
