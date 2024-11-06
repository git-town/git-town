Feature: display all executed Git commands

  Scenario: Git Town command ran successfully
    Given a Git repo with origin
    And I ran "git-town sync"
    When I run "git-town status --verbose"
    Then Git Town runs the commands
      | BRANCH | TYPE    | COMMAND                            |
      |        | backend | git version                        |
      |        | backend | git rev-parse --show-toplevel      |
      |        | backend | git config -lz --includes --global |
      |        | backend | git config -lz --includes --local  |
    And Git Town prints:
      """
      Ran 4 shell commands.
      """
