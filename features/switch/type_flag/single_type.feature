@messyoutput
Feature: switch branches of a single type

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
    When I run "git-town switch --type=observed" and enter into the dialogs:
      | DIALOG        | KEYS       |
      | switch-branch | down enter |
    Then Git Town runs the commands
      | BRANCH     | COMMAND                 |
      | observed-2 | git checkout observed-1 |

  Scenario: short form
    When I run "git-town switch -to" and enter into the dialogs:
      | DIALOG        | KEYS       |
      | switch-branch | down enter |
    Then Git Town runs the commands
      | BRANCH     | COMMAND                 |
      | observed-2 | git checkout observed-1 |
  #
  # NOTE: Cannot test undo here.
  # The Git Town command under test has not created an undoable runstate.
  # Executing "git town undo" would undo the Git Town command executed during setup.
