Feature: offline mode

  Background:
    Given a Git repo with origin
    And offline mode is enabled
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | old  | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
      | old    | local, origin | old commit  |
    And the current branch is "old"
    When I run "git-town rename new"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                   |
      | old    | git branch --move old new |
      |        | git checkout new          |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
      | new    | local         | old commit  |
      | old    | origin        | old commit  |
    And this lineage exists now
      | BRANCH | PARENT |
      | new    | main   |

  @debug
  @this
  Scenario: undo
    When I run "git-town undo -v"
    Then Git Town runs the commands
      | BRANCH | COMMAND          |
      | new    | git checkout old |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
      | new    | local         | old commit  |
      | old    | local, origin | old commit  |
    And these branches exist now
      | REPOSITORY | BRANCHES       |
      | local      | main, new, old |
      | origin     | main, old      |
