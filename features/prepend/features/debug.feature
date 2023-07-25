Feature: display debug statistics

  Background:
    Given the current branch is a feature branch "old"
    And the commits
      | BRANCH | LOCATION      | MESSAGE    |
      | old    | local, origin | old commit |

  Scenario: result
    When I run "git-town prepend parent --debug"
    Then it runs the commands
      | BRANCH | TYPE     | COMMAND                                       |
      |        | backend  | git version                                   |
      |        | backend  | git config -lz --local                        |
      |        | backend  | git config -lz --global                       |
      |        | backend  | git rev-parse --show-toplevel                 |
      |        | backend  | git remote                                    |
      |        | backend  | git status                                    |
      |        | backend  | git rev-parse --abbrev-ref HEAD               |
      | old    | frontend | git fetch --prune --tags                      |
      |        | backend  | git branch -vva                               |
      |        | backend  | git rev-parse --verify --abbrev-ref @{-1}     |
      |        | backend  | git status --porcelain --ignore-submodules    |
      | old    | frontend | git checkout main                             |
      |        | backend  | git rev-parse HEAD                            |
      | main   | frontend | git rebase origin/main                        |
      |        | backend  | git rev-list --left-right main...origin/main  |
      | main   | frontend | git checkout old                              |
      |        | backend  | git rev-parse HEAD                            |
      | old    | frontend | git merge --no-edit origin/old                |
      |        | backend  | git rev-parse HEAD                            |
      | old    | frontend | git merge --no-edit main                      |
      |        | backend  | git rev-list --left-right old...origin/old    |
      | old    | frontend | git branch parent main                        |
      |        | backend  | git config git-town-branch.parent.parent main |
      |        | backend  | git config git-town-branch.old.parent parent  |
      | old    | frontend | git checkout parent                           |
      |        | backend  | git branch                                    |
      |        | backend  | git branch                                    |
      |        | backend  | git rev-parse --verify --abbrev-ref @{-1}     |
    And it prints:
      """
      Ran 28 shell commands.
      """
    And the current branch is now "parent"

  Scenario: undo
    Given I ran "git-town prepend parent"
    When I run "git-town undo --debug"
    Then it runs the commands
      | BRANCH | TYPE     | COMMAND                                          |
      |        | backend  | git version                                      |
      |        | backend  | git config -lz --local                           |
      |        | backend  | git config -lz --global                          |
      |        | backend  | git rev-parse --show-toplevel                    |
      |        | backend  | git branch -vva                                  |
      | parent | frontend | git checkout old                                 |
      |        | backend  | git config git-town-branch.old.parent main       |
      |        | backend  | git config --unset git-town-branch.parent.parent |
      |        | backend  | git rev-parse parent                             |
      |        | backend  | git log main..parent                             |
      | old    | frontend | git branch -D parent                             |
      |        | frontend | git checkout main                                |
      | main   | frontend | git checkout old                                 |
    And it prints:
      """
      Ran 13 shell commands.
      """
    And the current branch is now "old"
