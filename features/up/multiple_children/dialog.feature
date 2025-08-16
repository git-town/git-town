@messyoutput
Feature: move up when branch has multiple children

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE    | PARENT   | LOCATIONS     |
      | branch-1  | feature | main     | local, origin |
      | branch-1a | feature | branch-1 | local, origin |
      | branch-1b | feature | branch-1 | local, origin |
    And the current branch is "branch-1"

  Scenario: select the first child
    When I run "git-town up" and enter into the dialogs:
      | DIALOG       | KEYS  |
      | child-branch | enter |
    Then Git Town runs the commands
      | BRANCH   | COMMAND                |
      | branch-1 | git checkout branch-1a |
    And Git Town prints:
      """
        main
          branch-1
      *     branch-1a
            branch-1b
      """

  Scenario: selecting the second child
    When I run "git-town up" and enter into the dialogs:
      | DIALOG       | KEYS       |
      | child-branch | down enter |
    Then Git Town runs the commands
      | BRANCH   | COMMAND                |
      | branch-1 | git checkout branch-1b |
    And Git Town prints:
      """
        main
          branch-1
            branch-1a
      *     branch-1b
      """

  Scenario: aborting the dialog
    When I run "git-town up" and enter into the dialogs:
      | DIALOG       | KEYS |
      | child-branch | q    |
    Then Git Town runs no commands
