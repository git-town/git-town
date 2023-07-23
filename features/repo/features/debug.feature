Feature: display debug statistics

  Scenario:
    Given the origin is "git@github.com:git-town/git-town.git"
    And tool "open" is installed
    When I run "git-town repo --debug"
    Then it runs the commands
      | BRANCH | TYPE     | COMMAND                                   |
      |        | backend  | git config -lz --local                    |
      |        | backend  | git config -lz --global                   |
      |        | backend  | git rev-parse                             |
      |        | backend  | git rev-parse --show-toplevel             |
      |        | backend  | git version                               |
      |        | backend  | git branch -a                             |
      |        | backend  | which wsl-open                            |
      |        | backend  | which garcon-url-handler                  |
      |        | backend  | which xdg-open                            |
      |        | backend  | which open                                |
      |        | backend  | git status                                |
      |        | backend  | git rev-parse --abbrev-ref HEAD           |
      | <none> | frontend | open https://github.com/git-town/git-town |
    And "open" launches a new pull request with this url in my browser:
      """
      https://github.com/git-town/git-town
      """
