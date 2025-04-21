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
    When I run "git-town sync --verbose"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | TYPE     | COMMAND                                          |
      |        | backend  | git version                                      |
      |        | backend  | git rev-parse --show-toplevel                    |
      |        | backend  | git config -lz --includes --global               |
      |        | backend  | git config -lz --includes --local                |
      |        | backend  | git -c core.abbrev=40 branch -vva --sort=refname |
      |        | backend  | git status -z --ignore-submodules                |
      |        | backend  | git rev-parse -q --verify MERGE_HEAD             |
      |        | backend  | git rev-parse -q --verify REBASE_HEAD            |
      |        | backend  | git remote                                       |
      | old    | frontend | git fetch --prune --tags                         |
      |        | backend  | git stash list                                   |
      |        | backend  | git -c core.abbrev=40 branch -vva --sort=refname |
      |        | backend  | git rev-parse --verify --abbrev-ref @{-1}        |
      |        | backend  | git remote get-url origin                        |
      |        | backend  | git log main..old --format=%s --reverse          |
      | old    | frontend | git checkout main                                |
      |        | backend  | git config --unset git-town-branch.old.parent    |
      | main   | frontend | git branch -D old                                |
      |        | backend  | git show-ref --verify --quiet refs/heads/old     |
      |        | backend  | git show-ref --verify --quiet refs/heads/active  |
      | main   | frontend | git checkout active                              |
      |        | backend  | git -c core.abbrev=40 branch -vva --sort=refname |
      |        | backend  | git config -lz --includes --global               |
      |        | backend  | git config -lz --includes --local                |
      |        | backend  | git stash list                                   |
    And Git Town prints:
      """
      Ran 25 shell commands.
      """
    And the branches are now
      | REPOSITORY    | BRANCHES     |
      | local, origin | main, active |
    And this lineage exists now
      | BRANCH | PARENT |
      | active | main   |

  Scenario: undo
    When I run "git-town undo --verbose"
    Then Git Town runs the commands
      | BRANCH | TYPE     | COMMAND                                          |
      |        | backend  | git version                                      |
      |        | backend  | git rev-parse --show-toplevel                    |
      |        | backend  | git config -lz --includes --global               |
      |        | backend  | git config -lz --includes --local                |
      |        | backend  | git status -z --ignore-submodules                |
      |        | backend  | git rev-parse -q --verify MERGE_HEAD             |
      |        | backend  | git rev-parse -q --verify REBASE_HEAD            |
      |        | backend  | git stash list                                   |
      |        | backend  | git -c core.abbrev=40 branch -vva --sort=refname |
      |        | backend  | git remote get-url origin                        |
      |        | backend  | git rev-parse --verify --abbrev-ref @{-1}        |
      |        | backend  | git remote get-url origin                        |
      | active | frontend | git branch old {{ sha 'initial commit' }}        |
      |        | backend  | git show-ref --quiet refs/heads/old              |
      | active | frontend | git checkout old                                 |
      |        | backend  | git config git-town-branch.old.parent main       |
    And Git Town prints:
      """
      Ran 16 shell commands.
      """
    And the initial branches and lineage exist now
