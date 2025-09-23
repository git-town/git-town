@messyoutput
Feature: switch branches using multiple types

  Background:
    Given a Git repo with origin
    And the branches
      | NAME         | TYPE         | PARENT | LOCATIONS     |
      | contribution | contribution |        | local, origin |
      | feature      | feature      | main   | local         |
      | observed-1   | observed     |        | local, origin |
      | observed-2   | observed     |        | local, origin |
      | parked       | parked       | main   | local         |
      | perennial    | perennial    |        | local         |
      | prototype    | prototype    | main   | local         |
    And the current branch is "observed-2"

  Scenario: long form
    When I run "git-town switch --type=observed+prototype" and enter into the dialogs:
      | DIALOG        | KEYS       |
      | switch-branch | down enter |
    Then Git Town runs the commands
      | BRANCH     | COMMAND                |
      | observed-2 | git checkout prototype |

  Scenario: short form
    When I run "git-town switch -to+pr" and enter into the dialogs:
      | DIALOG        | KEYS       |
      | switch-branch | down enter |
    Then Git Town runs the commands
      | BRANCH     | COMMAND                |
      | observed-2 | git checkout prototype |

  Scenario: undo
    Given I ran "git-town switch -to+pr" and enter into the dialogs:
      | DIALOG        | KEYS       |
      | switch-branch | down enter |
    When I run "git-town undo"
    Then Git Town runs no commands
    And the initial branches and lineage exist now
