Feature: update the parent of a feature branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | parent | feature | main   | local, origin |
      | child  | feature | parent | local, origin |
    And the current branch is "child"
    When I run "git-town set-parent main"

  Scenario: result
    Then Git Town runs no commands
    And this lineage exists now
      | BRANCH | PARENT |
      | child  | main   |
      | parent | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the initial commits exist now
    And the initial branches and lineage exist now
