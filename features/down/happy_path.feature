Feature: move down one position in the current stack

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
      | beta  | feature | alpha  | local, origin |
    And the current branch is "beta"
    When I run "git-town down"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND            |
      | beta   | git checkout alpha |
    And Git Town prints:
      """
        main
      *   alpha
            beta
      """
