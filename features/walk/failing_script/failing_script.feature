@skipwindows
Feature: if the given script returns an error, continue re-runs the failing script

  Scenario: action
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | branch-1 | feature | main   | local     |
      | branch-2 | feature | main   | local     |
    And the current branch is "branch-2"
    When I run "git-town walk --all test"
    Then Git Town runs the commands
      | BRANCH   | COMMAND               |
      | branch-2 | git checkout branch-1 |
      | branch-1 | test                  |
    And Git Town prints the error:
      """
      To continue after having resolved conflicts, run "git town continue".
      To go back to where you started, run "git town undo".
      """
    When I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH   | COMMAND |
      | branch-1 | test    |
    When I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH   | COMMAND |
      | branch-1 | test    |
