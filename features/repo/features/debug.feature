Feature: display debug statistics

  Scenario:
    Given the origin is "git@github.com:git-town/git-town.git"
    And tool "open" is installed
    When I run "git-town repo --debug"
    Then it runs the debug commands
      | git config -lz --local          |
      | git config -lz --global         |
      | git rev-parse                   |
      | git rev-parse --show-toplevel   |
      | git version                     |
      | git branch -a                   |
      | which wsl-open                  |
      | which garcon-url-handler        |
      | which xdg-open                  |
      | which open                      |
      | git status                      |
      | git rev-parse --abbrev-ref HEAD |
    And "open" launches a new pull request with this url in my browser:
      """
      https://github.com/git-town/git-town
      """
