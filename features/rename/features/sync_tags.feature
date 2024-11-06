Feature: don't sync tags while renaming branches

  Background:
    Given a Git repo with origin
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | old  | feature | main   | local, origin |
    And the tags
      | NAME       | LOCATION |
      | local-tag  | local    |
      | origin-tag | origin   |
    And the current branch is "old"
    And Git Town setting "sync-tags" is "false"
    When I run "git-town rename new"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                     |
      | old    | git fetch --prune --no-tags |
      |        | git branch --move old new   |
      |        | git checkout new            |
      | new    | git push -u origin new      |
      |        | git push origin :old        |
    And the initial tags exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                   |
      | new    | git branch old {{ sha 'initial commit' }} |
      |        | git push -u origin old                    |
      |        | git checkout old                          |
      | old    | git branch -D new                         |
      |        | git push origin :new                      |
    And the initial commits exist now
    And the initial lineage exists now
    And the initial tags exist now
