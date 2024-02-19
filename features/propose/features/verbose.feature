@skipWindows
Feature: display all executed Git commands

  Scenario: verbose mode enabled
    Given tool "open" is installed
    And the current branch is a feature branch "feature"
    And the origin is "git@github.com:git-town/git-town.git"
    When I run "git-town propose --verbose"
    Then it runs the commands
      | BRANCH  | TYPE     | COMMAND                                                            |
      |         | backend  | git version                                                        |
      |         | backend  | git config -lz --global                                            |
      |         | backend  | git config -lz --local                                             |
      |         | backend  | git rev-parse --show-toplevel                                      |
      |         | backend  | git stash list                                                     |
      |         | backend  | git status --long --ignore-submodules                              |
      |         | backend  | git branch -vva                                                    |
      |         | backend  | git remote                                                         |
      | feature | frontend | git fetch --prune --tags                                           |
      |         | backend  | git branch -vva                                                    |
      |         | backend  | git rev-parse --verify --abbrev-ref @{-1}                          |
      | feature | frontend | git checkout main                                                  |
      | main    | frontend | git rebase origin/main                                             |
      |         | backend  | git rev-list --left-right main...origin/main                       |
      | main    | frontend | git checkout feature                                               |
      | feature | frontend | git merge --no-edit origin/feature                                 |
      |         | frontend | git merge --no-edit main                                           |
      |         | backend  | git rev-list --left-right feature...origin/feature                 |
      |         | backend  | git show-ref --verify --quiet refs/heads/main                      |
      |         | backend  | which wsl-open                                                     |
      |         | backend  | which garcon-url-handler                                           |
      |         | backend  | which xdg-open                                                     |
      |         | backend  | which open                                                         |
      | <none>  | frontend | open https://github.com/git-town/git-town/compare/feature?expand=1 |
      |         | backend  | git branch -vva                                                    |
      |         | backend  | git config -lz --global                                            |
      |         | backend  | git config -lz --local                                             |
      |         | backend  | git stash list                                                     |
    And it prints:
      """
      Ran 27 shell commands.
      """
    And "open" launches a new proposal with this url in my browser:
      """
      https://github.com/git-town/git-town/compare/feature?expand=1
      """
