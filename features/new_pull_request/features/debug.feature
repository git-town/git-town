@skipWindows
Feature: display debug statistics

  Scenario: debug mode enabled
    Given tool "open" is installed
    And the current branch is a feature branch "feature"
    And the origin is "git@github.com:git-town/git-town.git"
    When I run "git-town new-pull-request --debug"
    Then it runs the debug commands
      | git config -lz --local                             |
      | git config -lz --global                            |
      | git rev-parse                                      |
      | git rev-parse --show-toplevel                      |
      | git version                                        |
      | git branch -a                                      |
      | git remote                                         |
      | git status                                         |
      | git rev-parse --abbrev-ref HEAD                    |
      | git branch -r                                      |
      | git rev-parse --verify --abbrev-ref @{-1}          |
      | git status --porcelain --ignore-submodules         |
      | git rev-parse HEAD                                 |
      | git rev-list --left-right main...origin/main       |
      | git rev-parse HEAD                                 |
      | git rev-parse HEAD                                 |
      | git rev-list --left-right feature...origin/feature |
      | git branch                                         |
      | git branch                                         |
      | git rev-parse --verify --abbrev-ref @{-1}          |
      | which wsl-open                                     |
      | which garcon-url-handler                           |
      | which xdg-open                                     |
      | which open                                         |
    And "open" launches a new pull request with this url in my browser:
      """
      https://github.com/git-town/git-town/compare/feature?expand=1
      """
