@skipWindows
Feature: display all executed Git commands

  Scenario:
    Given the origin is "git@github.com:git-town/git-town.git"
    And tool "open" is installed
    When I run "git-town repo --verbose"
    Then it runs the commands
      | BRANCH | TYPE     | COMMAND                                   |
      |        | backend  | git version                               |
      |        | backend  | git config -lz --global                   |
      |        | backend  | git config -lz --local                    |
      |        | backend  | git rev-parse --show-toplevel             |
      |        | backend  | which wsl-open                            |
      |        | backend  | which garcon-url-handler                  |
      |        | backend  | which xdg-open                            |
      |        | backend  | which open                                |
      |        | backend  | git status --long --ignore-submodules     |
      |        | backend  | git rev-parse --abbrev-ref HEAD           |
      | <none> | frontend | open https://github.com/git-town/git-town |
    And it prints:
      """
      Ran 11 shell commands.
      """
    And "open" launches a new proposal with this url in my browser:
      """
      https://github.com/git-town/git-town
      """
