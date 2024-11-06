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
      |        | backend  | git config -lz --includes --global        |
      |        | backend  | git config -lz --includes --local         |
      |        | backend  | which wsl-open                            |
      |        | backend  | which garcon-url-handler                  |
      |        | backend  | which xdg-open                            |
      |        | backend  | which open                                |
      |        | backend  | git rev-parse --abbrev-ref HEAD           |
      | <none> | frontend | open https://github.com/git-town/git-town |
    And Git Town prints:
      """
      Ran 10 shell commands.
      """
    And "open" launches a new proposal with this url in my browser:
      """
      https://github.com/git-town/git-town
      """
