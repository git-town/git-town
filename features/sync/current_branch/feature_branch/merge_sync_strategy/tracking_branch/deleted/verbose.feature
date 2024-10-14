Feature: display all executed Git commands

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | active | feature | main   | local, origin |
      | old    | feature | main   | local, origin |
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
      |        | backend  | git rev-parse --show-toplevel                 |
      |        | backend  | git config -lz --includes --global            |
      |        | backend  | git config -lz --includes --local             |
      |        | backend  | git branch -vva --sort=refname                |
      |        | backend  | git status --long --ignore-submodules         |
      |        | backend  | git remote                                    |
      | old    | frontend | git fetch --prune --tags                      |
      |        | backend  | git stash list                                |
      |        | backend  | git branch -vva --sort=refname                |
      |        | backend  | git rev-parse --verify --abbrev-ref @{-1}     |
      |        | backend  | git remote get-url origin                     |
      |        | backend  | git log main..old --format=%h                 |
      | old    | frontend | git checkout main                             |
      | main   | frontend | git rebase origin/main --no-update-refs       |
      |        | backend  | git rev-list --left-right main...origin/main  |
      |        | backend  | git config --unset git-town-branch.old.parent |
      | main   | frontend | git branch -D old                             |
      |        | backend  | git show-ref --verify --quiet refs/heads/old  |
      |        | backend  | git show-ref --verify --quiet refs/heads/main |
      |        | backend  | git branch -vva --sort=refname                |
      |        | backend  | git config -lz --includes --global            |
      |        | backend  | git config -lz --includes --local             |
      |        | backend  | git stash list                                |
    And it prints:
      """
      Ran 24 shell commands.
      """
    And the current branch is now "main"
    And the branches are now
      | REPOSITORY    | BRANCHES     |
      | local, origin | main, active |
    And this lineage exists now
      | BRANCH | PARENT |
      | active | main   |

  Scenario: undo
    Given I ran "git-town sync"
    When I run "git-town undo --verbose"
    Then it runs the commands
      | BRANCH | TYPE     | COMMAND                                    |
      |        | backend  | git version                                |
      |        | backend  | git rev-parse --show-toplevel              |
      |        | backend  | git config -lz --includes --global         |
      |        | backend  | git config -lz --includes --local          |
      |        | backend  | git status --long --ignore-submodules      |
      |        | backend  | git stash list                             |
      |        | backend  | git branch -vva --sort=refname             |
      |        | backend  | git remote get-url origin                  |
      |        | backend  | git rev-parse --verify --abbrev-ref @{-1}  |
      |        | backend  | git remote get-url origin                  |
      | main   | frontend | git branch old {{ sha 'initial commit' }}  |
      |        | backend  | git show-ref --quiet refs/heads/old        |
      | main   | frontend | git checkout old                           |
      |        | backend  | git config git-town-branch.old.parent main |
    And it prints:
      """
      Ran 14 shell commands.
      """
    And the current branch is now "old"
    And the initial branches and lineage exist now
