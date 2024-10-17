Feature: display all executed Git commands

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | branch-1 | feature | main   | local, origin |
      | branch-2 | feature | main   | local, origin |
    And Git Town setting "sync-feature-strategy" is "rebase"
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         |
      | branch-1 | local, origin | branch-1 commit |
    And origin deletes the "branch-2" branch
    And the current branch is "branch-2"

  Scenario: result
    When I run "git-town sync --verbose"
    Then it runs the commands
      | BRANCH   | TYPE     | COMMAND                                            |
      |          | backend  | git version                                        |
      |          | backend  | git rev-parse --show-toplevel                      |
      |          | backend  | git config -lz --includes --global                 |
      |          | backend  | git config -lz --includes --local                  |
      |          | backend  | git branch -vva --sort=refname                     |
      |          | backend  | git status --long --ignore-submodules              |
      |          | backend  | git remote                                         |
      | branch-2 | frontend | git fetch --prune --tags                           |
      |          | backend  | git stash list                                     |
      |          | backend  | git branch -vva --sort=refname                     |
      |          | backend  | git rev-parse --verify --abbrev-ref @{-1}          |
      |          | backend  | git remote get-url origin                          |
      |          | backend  | git log main..branch-2 --format=%h                 |
      | branch-2 | frontend | git checkout main                                  |
      | main     | frontend | git rebase origin/main --no-update-refs            |
      |          | backend  | git rev-list --left-right main...origin/main       |
      |          | backend  | git config --unset git-town-branch.branch-2.parent |
      | main     | frontend | git branch -D branch-2                             |
      |          | backend  | git show-ref --verify --quiet refs/heads/branch-2  |
      |          | backend  | git show-ref --verify --quiet refs/heads/main      |
      |          | backend  | git branch -vva --sort=refname                     |
      |          | backend  | git config -lz --includes --global                 |
      |          | backend  | git config -lz --includes --local                  |
      |          | backend  | git stash list                                     |
    And it prints:
      """
      Ran 24 shell commands.
      """
    And the current branch is now "main"
    And the branches are now
      | REPOSITORY    | BRANCHES       |
      | local, origin | main, branch-1 |
    And this lineage exists now
      | BRANCH   | PARENT |
      | branch-1 | main   |

  Scenario: undo
    Given I ran "git-town sync"
    When I run "git-town undo --verbose"
    Then it runs the commands
      | BRANCH | TYPE     | COMMAND                                         |
      |        | backend  | git version                                     |
      |        | backend  | git rev-parse --show-toplevel                   |
      |        | backend  | git config -lz --includes --global              |
      |        | backend  | git config -lz --includes --local               |
      |        | backend  | git status --long --ignore-submodules           |
      |        | backend  | git stash list                                  |
      |        | backend  | git branch -vva --sort=refname                  |
      |        | backend  | git remote get-url origin                       |
      |        | backend  | git rev-parse --verify --abbrev-ref @{-1}       |
      |        | backend  | git remote get-url origin                       |
      | main   | frontend | git branch branch-2 {{ sha 'initial commit' }}  |
      |        | backend  | git show-ref --quiet refs/heads/branch-2        |
      | main   | frontend | git checkout branch-2                           |
      |        | backend  | git config git-town-branch.branch-2.parent main |
    And it prints:
      """
      Ran 14 shell commands.
      """
    And the current branch is now "branch-2"
    And the initial branches and lineage exist now
