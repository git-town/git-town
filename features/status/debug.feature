Feature: display debug statistics

  Scenario: Git Town command ran successfully
    Given I ran "git-town sync"
    When I run "git-town status --debug"
    Then it runs the commands
      | BRANCH | TYPE    | COMMAND                       |
      |        | backend | git version                   |
      |        | backend | git config -lz --local        |
      |        | backend | git config -lz --global       |
      |        | backend | git rev-parse --show-toplevel |
    And it prints:
      """
      Ran 4 shell commands.
      """
