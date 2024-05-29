Feature: display all executed Git commands

  Background:
    Given the current branch is a feature branch "old"
    And the commits
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
      | old    | local, origin | old commit  |

  Scenario: result
    When I run "git-town rename-branch new --verbose"
    Then it runs the commands
      | BRANCH | TYPE     | COMMAND                                       |
      |        | backend  | git version                                   |
      |        | backend  | git config -lz --global                       |
      |        | backend  | git config -lz --local                        |
      |        | backend  | git rev-parse --show-toplevel                 |
      |        | backend  | git rev-parse --verify --abbrev-ref @{-1}     |
      |        | backend  | git status --long --ignore-submodules         |
      |        | backend  | git remote                                    |
      |        | backend  | git rev-parse --abbrev-ref HEAD               |
      | old    | frontend | git fetch --prune --tags                      |
      |        | backend  | git stash list                                |
      |        | backend  | git branch -vva --sort=refname                |
      | old    | frontend | git branch new old                            |
      |        | frontend | git checkout new                              |
      |        | backend  | git config --unset git-town-branch.old.parent |
      |        | backend  | git config git-town-branch.new.parent main    |
      | new    | frontend | git push -u origin new                        |
      |        | frontend | git push origin :old                          |
      |        | frontend | git branch -D old                             |
      |        | backend  | git show-ref --verify --quiet refs/heads/main |
      |        | backend  | git checkout main                             |
      |        | backend  | git checkout new                              |
      |        | backend  | git branch -vva --sort=refname                |
      |        | backend  | git config -lz --global                       |
      |        | backend  | git config -lz --local                        |
      |        | backend  | git stash list                                |
    And it prints:
      """
      Ran 25 shell commands.
      """
    And the current branch is now "new"

  Scenario: undo
    Given I ran "git-town rename-branch new"
    When I run "git-town undo --verbose"
    Then it runs the commands
      | BRANCH | TYPE     | COMMAND                                       |
      |        | backend  | git version                                   |
      |        | backend  | git config -lz --global                       |
      |        | backend  | git config -lz --local                        |
      |        | backend  | git rev-parse --show-toplevel                 |
      |        | backend  | git status --long --ignore-submodules         |
      |        | backend  | git stash list                                |
      |        | backend  | git branch -vva --sort=refname                |
      |        | backend  | git rev-parse --verify --abbrev-ref @{-1}     |
      |        | backend  | git remote get-url origin                     |
      | new    | frontend | git branch old {{ sha 'old commit' }}         |
      |        | frontend | git push -u origin old                        |
      |        | frontend | git push origin :new                          |
      |        | frontend | git checkout old                              |
      | old    | frontend | git branch -D new                             |
      |        | backend  | git config --unset git-town-branch.new.parent |
      |        | backend  | git config git-town-branch.old.parent main    |
    And it prints:
      """
      Ran 16 shell commands.
      """
    And the current branch is now "old"
