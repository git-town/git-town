Feature: offline mode

  Background:
    Given a Git repo with origin
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | old  | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
      | old    | local, origin | old commit  |
    And offline mode is enabled
    And the current branch is "old"
    When I run "git-town rename new"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                   |
      | old    | git branch --move old new |
      |        | git checkout new          |
    And this lineage exists now
      """
      main
        new
      """
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
      | new    | local         | old commit  |
      | old    | origin        | old commit  |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                   |
      | new    | git branch --move new old |
      |        | git checkout old          |
    And the initial branches and lineage exist now
    And the initial commits exist now
