@skipWindows
Feature: display all executed Git commands

  Scenario:
    Given a Git repo with origin
    And the origin is "git@github.com:git-town/git-town.git"
    And tool "open" is installed
    When I run "git-town repo --verbose"
    Then Git Town runs the commands
      | BRANCH | TYPE     | COMMAND                                   |
      |        | backend  | git version                               |
      |        | backend  | git rev-parse --show-toplevel             |
      |        | backend  | git config -lz --global                   |
      |        | backend  | git config -lz --local                    |
      |        | backend  | git config -lz                            |
      |        | backend  | git branch --show-current                 |
      | main   | frontend | open https://github.com/git-town/git-town |
    And Git Town prints:
      """
      Ran 7 shell commands.
      """
