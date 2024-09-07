Feature: switch branches using multiple types

  Background:
    Given a Git repo with origin
    And the branches
      | NAME         | TYPE         | PARENT | LOCATIONS |
      | contribution | contribution |        | local     |
      | feature      | feature      | main   | local     |
      | observed-1   | observed     |        | local     |
      | observed-2   | observed     |        | local     |
      | parked       | parked       | main   | local     |
      | perennial    | perennial    |        | local     |
      | prototype    | prototype    | main   | local     |
    And the current branch is "observed-2"
    When I run "git-town switch --type=observed+prototype" and enter into the dialogs:
      | KEYS       |
      | down enter |

  Scenario: switching to another branch
    Then it runs the commands
      | BRANCH     | COMMAND                |
      | observed-2 | git checkout prototype |
    And the current branch is now "prototype"

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "prototype"
    And the initial branches and lineage exist
