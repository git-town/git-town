Feature: display all executed Git commands

  Scenario: Git Town command ran successfully
    Given a Git repo clone
    Given I ran "git-town sync"
    When I run "git-town status --verbose"
    Then it runs the commands
      | BRANCH | TYPE    | COMMAND                            |
      |        | backend | git version                        |
      |        | backend | git rev-parse --show-toplevel      |
      |        | backend | git config -lz --includes --global |
      |        | backend | git config -lz --includes --local  |
    And it prints:
      """
      Ran 4 shell commands.
      """
