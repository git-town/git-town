Feature: display debug statistics

  Scenario: Git Town command ran successfully
    Given I ran "git-town sync"
    When I run "git-town status --debug"
    Then it runs the commands
      | BRANCH | TYPE    | COMMAND                       |
      |        | backend | git config -lz --local        |
      |        | backend | git config -lz --global       |
      |        | backend | git rev-parse                 |
      |        | backend | git rev-parse --show-toplevel |
      |        | backend | git version                   |
    And it prints:
      """
      Ran 5 shell commands.
      """
