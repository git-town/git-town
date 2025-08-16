Feature: move up one position in the current stack

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
      | beta  | feature | alpha  | local, origin |
    And the current branch is "alpha"
    When I run "git-town up"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND           |
      | alpha  | git checkout beta |
    And Git Town prints:
      """
        main
          alpha
      *     beta
      """
