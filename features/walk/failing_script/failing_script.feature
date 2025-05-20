@skipWindows
Feature: handle errors in the given script

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | branch-1 | feature | main   | local     |
      | branch-2 | feature | main   | local     |
      | branch-3 | feature | main   | local     |
    And the current branch is "branch-2"
    When I run "git-town walk --all test"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND               |
      | branch-2 | git checkout branch-1 |
      | branch-1 | test                  |
    And Git Town prints the error:
      """
      To continue after having resolved conflicts, run "git town continue".
      To go back to where you started, run "git town undo".
      """

  Scenario: continue re-runs the failed script
    When I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH   | COMMAND |
      | branch-1 | test    |
    And Git Town prints the error:
      """
      To continue after having resolved conflicts, run "git town continue".
      To go back to where you started, run "git town undo".
      """

  Scenario: skip runs the given script on the next branch
    When I run "git-town skip"
    Then Git Town runs the commands
      | BRANCH   | COMMAND               |
      | branch-1 | git checkout branch-2 |
      | branch-2 | test                  |
    And Git Town prints the error:
      """
      To continue after having resolved conflicts, run "git town continue".
      To go back to where you started, run "git town undo".
      """
    When I run "git-town skip"
    Then Git Town runs the commands
      | BRANCH   | COMMAND               |
      | branch-2 | git checkout branch-3 |
      | branch-3 | test                  |
    And Git Town prints the error:
      """
      To continue after having resolved conflicts, run "git town continue".
      To go back to where you started, run "git town undo".
      """
    When I run "git-town skip"
    Then Git Town runs the commands
      | BRANCH   | COMMAND               |
      | branch-3 | git checkout branch-2 |
    And Git Town prints:
      """
      Branch walk done.
      """
