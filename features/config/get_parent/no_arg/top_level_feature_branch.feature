Feature: display the parent of a top-level feature branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the current branch is "feature"
    When I run "git-town config get-parent"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
      main
      """
