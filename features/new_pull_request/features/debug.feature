@skipWindows
Feature: display debug statistics

  Scenario: debug mode enabled
    Given tool "open" is installed
    And the current branch is a feature branch "feature"
    And the origin is "git@github.com:git-town/git-town.git"
    When I run "git-town new-pull-request --debug"
    Then it runs the commands
      | BRANCH  | TYPE     | COMMAND                                                            |
      |         | backend  | git version                                                        |
      |         | backend  | git config -lz --local                                             |
      |         | backend  | git config -lz --global                                            |
      |         | backend  | git rev-parse --show-toplevel                                      |
      |         | backend  | git branch -vva                                                    |
      |         | backend  | git remote                                                         |
      | feature | frontend | git fetch --prune --tags                                           |
      |         | backend  | git rev-parse --verify --abbrev-ref @{-1}                          |
      |         | backend  | git status --porcelain --ignore-submodules                         |
      | feature | frontend | git checkout main                                                  |
      |         | backend  | git rev-parse HEAD                                                 |
      | main    | frontend | git rebase origin/main                                             |
      |         | backend  | git rev-list --left-right main...origin/main                       |
      | main    | frontend | git checkout feature                                               |
      |         | backend  | git rev-parse HEAD                                                 |
      | feature | frontend | git merge --no-edit origin/feature                                 |
      |         | backend  | git rev-parse HEAD                                                 |
      | feature | frontend | git merge --no-edit main                                           |
      |         | backend  | git rev-list --left-right feature...origin/feature                 |
      |         | backend  | git branch                                                         |
      |         | backend  | git branch                                                         |
      |         | backend  | git rev-parse --verify --abbrev-ref @{-1}                          |
      |         | backend  | which wsl-open                                                     |
      |         | backend  | which garcon-url-handler                                           |
      |         | backend  | which xdg-open                                                     |
      |         | backend  | which open                                                         |
      | <none>  | frontend | open https://github.com/git-town/git-town/compare/feature?expand=1 |
    And it prints:
      """
      Ran 27 shell commands.
      """
    And "open" launches a new pull request with this url in my browser:
      """
      https://github.com/git-town/git-town/compare/feature?expand=1
      """
