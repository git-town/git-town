@skipWindows
Feature: no TTY

  Scenario: main branch
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    And the current branch is "feature"
    When I run "git-town diff-parent" in a non-TTY shell
    Then Git Town runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git diff --merge-base main feature |
