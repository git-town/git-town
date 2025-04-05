Feature: reproduce bug

  Background:
    Given a Git repo with origin
    And I ran "git checkout -b test"
    And the commits
      | BRANCH | LOCATION | MESSAGE  | FILE NAME |
      | test   | local    | commit 1 | file_1    |
      | test   | local    | commit 2 | file_2    |
    And local Git setting "git-town.sync-feature-strategy" is "rebase"
    When I run "git-town set-parent" and enter into the dialog:
      | DIALOG                | KEYS  |
      | parent branch of test | enter |

  Scenario: result
    Then Git Town prints:
      """
      Selected parent branch for "test": main
      """
    And Git Town runs the commands
      | BRANCH | COMMAND                          |
      | test   | git rebase main --no-update-refs |
    And these commits exist now
      | BRANCH | LOCATION | MESSAGE  |
      | test   | local    | commit 1 |
      |        |          | commit 2 |
    And this lineage exists now
      | BRANCH | PARENT |
      | test   | main   |

  Scenario: undo
    When I run "git-town undo"
    And Git Town runs no commands
    And these commits exist now
      | BRANCH | LOCATION | MESSAGE  | FILE NAME |
      | test   | local    | commit 1 | file_1    |
      |        |          | commit 2 | file_2    |
    And the branches are now
      | REPOSITORY | BRANCHES   |
      | local      | main, test |
      | origin     | main       |
    And this lineage exists now
      | BRANCH | PARENT |
