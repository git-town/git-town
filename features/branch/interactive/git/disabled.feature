@skipWindows
Feature: interactivity disabled, no main branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
      | beta  | feature | alpha  | local, origin |
    And Git Town is not configured
    And Git setting "git-town.interactive" is "false"
    And the current branch is "beta"
    When I run "git-town branch"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
        main
          alpha
      *     beta
      """
