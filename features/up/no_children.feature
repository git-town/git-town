Feature: switching when the branch has no children

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
    And the current branch is "alpha"
    When I run "git-town up"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      branch "alpha" has no children
      """
