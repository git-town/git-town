Feature: display all executed Git commands

  Background:
    Given the feature branches "active" and "old"
    And the commits
      | BRANCH | LOCATION      | MESSAGE       |
      | active | local, origin | active commit |
    And origin deletes the "old" branch
    And the current branch is "old"

  Scenario: result
    When I run "git-town sync --verbose"
    Then it runs the commands
      | BRANCH | TYPE     | COMMAND                                       |
      |        | backend  | git version                                   |
      |        | backend  | git config -lz --global                       |
      |        | backend  | git config -lz --local                        |
      |        | backend  | git rev-parse --show-toplevel                 |
      |        | backend  | git status --long --ignore-submodules         |
      |        | backend  | git stash list                                |
      |        | backend  | git branch -vva                               |
      |        | backend  | git remote                                    |
      | old    | frontend | git fetch --prune --tags                      |
      |        | backend  | git branch -vva                               |
      |        | backend  | git rev-parse --verify --abbrev-ref @{-1}     |
      | old    | frontend | git checkout main                             |
      | main   | frontend | git rebase origin/main                        |
      |        | backend  | git rev-list --left-right main...origin/main  |
      | main   | frontend | git checkout old                              |
      | old    | frontend | git merge --no-edit main                      |
      |        | backend  | git diff main..old                            |
      | old    | frontend | git checkout main                             |
      |        | backend  | git log main..old                             |
      | main   | frontend | git branch -d old                             |
      |        | backend  | git config --unset git-town-branch.old.parent |
      |        | backend  | git show-ref --quiet refs/heads/old           |
      |        | backend  | git branch -vva                               |
      |        | backend  | git config -lz --global                       |
      |        | backend  | git config -lz --local                        |
      |        | backend  | git stash list                                |
    And it prints:
      """
      Ran 26 shell commands.
      """
    And the current branch is now "main"
    And the branches are now
      | REPOSITORY    | BRANCHES     |
      | local, origin | main, active |
    And this branch lineage exists now
      | BRANCH | PARENT |
      | active | main   |

  Scenario: undo
    Given I ran "git-town sync"
    When I run "git-town undo --verbose"
    Then it runs the commands
      | BRANCH | TYPE     | COMMAND                                    |
      |        | backend  | git version                                |
      |        | backend  | git config -lz --global                    |
      |        | backend  | git config -lz --local                     |
      |        | backend  | git rev-parse --show-toplevel              |
      |        | backend  | git status --long --ignore-submodules      |
      |        | backend  | git stash list                             |
      |        | backend  | git branch -vva                            |
      |        | backend  | git rev-parse --verify --abbrev-ref @{-1}  |
      |        | backend  | git remote get-url origin                  |
      | main   | frontend | git branch old {{ sha 'initial commit' }}  |
      |        | backend  | git show-ref --quiet refs/heads/old        |
      | main   | frontend | git checkout old                           |
      |        | backend  | git config git-town-branch.old.parent main |
    And it prints:
      """
      Ran 13 shell commands.
      """
    And the current branch is now "old"
    And the initial branches and lineage exist
