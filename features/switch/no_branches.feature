Feature: switch branches

  Scenario: no branches to switch to
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS |
      | alpha | feature | main   | local     |
    And the current branch is "alpha"
    When I run "git-town switch zonk"
    Then it runs no commands
    And it prints the error:
      """
      no branches to switch to
      """
    And the current branch is still "alpha"
