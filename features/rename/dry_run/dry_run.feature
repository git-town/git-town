Feature: rename the current branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | old  | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
      | old    | local, origin | old commit  |
    And the current branch is "old"
    When I run "git-town rename new --dry-run"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                   |
      | old    | git fetch --prune --tags  |
      |        | git branch --move old new |
      |        | git checkout new          |
      | new    | git push -u origin new    |
      |        | git push origin :old      |
    And the initial branches and lineage exist now
    And the initial commits exist now
  #
  # Cannot test undo because dry-run now doesn't create a runstate.
