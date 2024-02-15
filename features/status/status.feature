@smoke
Feature: describe the status of the current/last Git Town command

  Scenario: Git Town command ran successfully
    Given I ran "git-town sync"
    When I run "git-town status"
    Then it prints:
      """
      The previous Git Town command (sync) finished successfully.
      You can run "git town undo" to go back to where you started.
      """

  Scenario: Git Town command in progress
    Given the current branch is a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main    | local    | conflicting main commit    | conflicting_file | main content    |
      | feature | local    | conflicting feature commit | conflicting_file | feature content |
    And I run "git-town sync"
    When I run "git-town status"
    Then it prints something like:
      """
      The last Git Town command \(sync\) hit a problem .*ms ago.
      You can run "git town undo" to go back to where you started.
      You can run "git town continue" to finish it.
      You can run "git town skip" to skip the currently failing operation.
      """

  Scenario: no runstate exists
    When I run "git-town status"
    Then it prints:
      """
      No status file found for this repository.
      """
