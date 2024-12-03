Feature: don't sync tags while deleting branches

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | current | feature | main   | local, origin |
    And the tags
      | NAME       | LOCATION |
      | local-tag  | local    |
      | origin-tag | origin   |
    And the current branch is "current"
    And Git Town setting "sync-tags" is "false"
    When I run "git-town delete"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                     |
      | current | git fetch --prune --no-tags |
      |         | git push origin :current    |
      |         | git checkout main           |
      | main    | git branch -D current       |
    And the initial tags exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                       |
      | main   | git branch current {{ sha 'initial commit' }} |
      |        | git push -u origin current                    |
      |        | git checkout current                          |
    And the initial commits exist now
    And the initial lineage exists now
    And the initial tags exist now
