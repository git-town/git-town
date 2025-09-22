Feature: display all executed Git commands

  Scenario: Git Town command ran successfully
    Given a Git repo with origin
    And I ran "git-town sync"
    When I run "git-town status --verbose"
    Then Git Town runs the commands
      | BRANCH | TYPE    | COMMAND                       |
      |        | backend | git version                   |
      |        | backend | git rev-parse --show-toplevel |
      |        | backend | git config -lz --global       |
      |        | backend | git config -lz --local        |
      |        | backend | git config -lz                |
    And Git Town prints:
      """
      Ran 5 shell commands.
      """
