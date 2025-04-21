@skipWindows
Feature: display all executed Git commands

  Scenario: verbose mode enabled
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And tool "open" is installed
    And the current branch is "feature"
    And the origin is "git@github.com:git-town/git-town.git"
    And a proposal for this branch does not exist
    When I run "git-town propose --verbose"
    Then Git Town runs the commands
      | BRANCH  | TYPE     | COMMAND                                                            |
      |         | backend  | git version                                                        |
      |         | backend  | git rev-parse --show-toplevel                                      |
      |         | backend  | git config -lz --includes --global                                 |
      |         | backend  | git config -lz --includes --local                                  |
      |         | backend  | git -c core.abbrev=40 branch -vva --sort=refname                   |
      |         | backend  | git status -z --ignore-submodules                                  |
      |         | backend  | git rev-parse -q --verify MERGE_HEAD                               |
      |         | backend  | git rev-parse -q --verify REBASE_HEAD                              |
      |         | backend  | git remote                                                         |
      | feature | frontend | git fetch --prune --tags                                           |
      |         | backend  | git stash list                                                     |
      |         | backend  | git -c core.abbrev=40 branch -vva --sort=refname                   |
      |         | backend  | git rev-parse --verify --abbrev-ref @{-1}                          |
      | (none)  | frontend | Looking for proposal online ... ok                                 |
      |         | backend  | git log main..feature --format=%s --reverse                        |
      | feature | frontend | git merge --no-edit --ff main                                      |
      |         | frontend | git merge --no-edit --ff origin/feature                            |
      |         | backend  | git show-ref --verify --quiet refs/heads/feature                   |
      |         | backend  | git rev-list --left-right feature...origin/feature                 |
      |         | backend  | git rev-parse --abbrev-ref --symbolic-full-name @{u}               |
      |         | backend  | git show-ref --verify --quiet refs/heads/main                      |
      |         | backend  | git checkout main                                                  |
      |         | backend  | git checkout feature                                               |
      |         | backend  | which wsl-open                                                     |
      |         | backend  | which garcon-url-handler                                           |
      |         | backend  | which xdg-open                                                     |
      |         | backend  | which open                                                         |
      | (none)  | frontend | open https://github.com/git-town/git-town/compare/feature?expand=1 |
      |         | backend  | git -c core.abbrev=40 branch -vva --sort=refname                   |
      |         | backend  | git config -lz --includes --global                                 |
      |         | backend  | git config -lz --includes --local                                  |
      |         | backend  | git stash list                                                     |
    And Git Town prints:
      """
      Ran 31 shell commands.
      """
