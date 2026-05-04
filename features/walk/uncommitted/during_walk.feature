@skipWindows
Feature: handle created uncommitted changes

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | branch-1 | feature | main   | local     |
      | branch-2 | feature | main   | local     |
      | branch-3 | feature | main   | local     |
    And the commits
      | BRANCH   | LOCATION | MESSAGE  | FILE NAME | FILE CONTENT |
      | branch-1 | local    | commit 1 | file      | content 1    |
      | branch-2 | local    | commit 2 | file      | content 2    |
      | branch-3 | local    | commit 3 | file      | content 3    |
    And the current branch is "branch-2"
    And tool format is installed
    When I run "git-town walk --all format"

  @this
  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND               |
      | branch-2 | git checkout branch-1 |
      | branch-1 | format                |
    And Git Town prints the error:
      """
      Uncommitted changes detected.
      """
    And Git Town prints the error:
      """
      To continue after having committed the changes, run "git town continue".
      To continue with the uncommitted changes on the next branch, run "git town skip".
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
