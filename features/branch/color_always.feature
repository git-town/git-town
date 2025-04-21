Feature: display the local branch hierarchy when I use color.ui=always

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
      | beta  | feature | alpha  | local, origin |
      | gamma | feature | beta   | local, origin |
    And the current branch is "beta"
    And local Git setting "color.ui" is "always"
    When I run "git-town branch"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
        main
          alpha
      *     beta
              gamma
      """
