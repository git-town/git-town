Feature: don't sync tags while renaming branches

  Background:
    Given a Git repo with origin
    And the branch
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | old  | feature | main   | local, origin |
    And the tags
      | NAME       | LOCATION |
      | local-tag  | local    |
      | origin-tag | origin   |
    And the current branch is "old"
    And Git Town setting "sync-tags" is "false"
    When I run "git-town rename-branch new"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                     |
      | old    | git fetch --prune --no-tags |
      |        | git branch new old          |
      |        | git checkout new            |
      | new    | git push -u origin new      |
      |        | git push origin :old        |
      |        | git branch -D old           |
    And the initial tags exist now

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                   |
      | new    | git branch old {{ sha 'initial commit' }} |
      |        | git push -u origin old                    |
      |        | git push origin :new                      |
      |        | git checkout old                          |
      | old    | git branch -D new                         |
    And the initial commits exist
    And the initial lineage exists
    And the initial tags exist now
