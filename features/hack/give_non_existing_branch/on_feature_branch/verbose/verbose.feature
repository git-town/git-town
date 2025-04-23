Feature: display all executed Git commands

  Background:
    Given a Git repo with origin
    And the commits
      | BRANCH | LOCATION | MESSAGE     |
      | main   | origin   | main commit |
    And the current branch is "main"
    When I run "git-town hack new --verbose"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | TYPE     | COMMAND                                           |
      |        | backend  | git version                                       |
      |        | backend  | git rev-parse --show-toplevel                     |
      |        | backend  | git config -lz --includes --global                |
      |        | backend  | git config -lz --includes --local                 |
      |        | backend  | git -c core.abbrev=40 branch -vva --sort=refname  |
      |        | backend  | git rev-parse --verify --abbrev-ref @{-1}         |
      |        | backend  | git status -z --ignore-submodules                 |
      |        | backend  | git rev-parse -q --verify MERGE_HEAD              |
      |        | backend  | git rev-parse --absolute-git-dir                  |
      |        | backend  | git remote                                        |
      | main   | frontend | git fetch --prune --tags                          |
      |        | backend  | git stash list                                    |
      |        | backend  | git -c core.abbrev=40 branch -vva --sort=refname  |
      |        | backend  | git remote get-url origin                         |
      | main   | frontend | git -c rebase.updateRefs=false rebase origin/main |
      |        | backend  | git rev-parse --verify -q refs/heads/main         |
      |        | backend  | git rev-list --left-right main...origin/main      |
      |        | backend  | git rev-parse --verify -q refs/heads/main         |
      | main   | frontend | git checkout -b new                               |
      |        | backend  | git rev-parse --verify -q refs/heads/main         |
      |        | backend  | git config git-town-branch.new.parent main        |
      |        | backend  | git rev-parse --verify -q refs/heads/main         |
      |        | backend  | git -c core.abbrev=40 branch -vva --sort=refname  |
      |        | backend  | git config -lz --includes --global                |
      |        | backend  | git config -lz --includes --local                 |
      |        | backend  | git stash list                                    |
    And Git Town prints:
      """
      Ran 26 shell commands.
      """

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
      |        | backend  | git rev-parse --absolute-git-dir                 |
      |        | backend  | git stash list                                   |
      |        | backend  | git -c core.abbrev=40 branch -vva --sort=refname |
      |        | backend  | git remote get-url origin                        |
      |        | backend  | git rev-parse --verify --abbrev-ref @{-1}        |
      |        | backend  | git remote get-url origin                        |
      | new    | frontend | git checkout main                                |
      |        | backend  | git rev-parse HEAD                               |
      | main   | frontend | git reset --hard {{ sha 'initial commit' }}      |
      |        | frontend | git branch -D new                                |
      |        | backend  | git config --unset git-town-branch.new.parent    |
    And Git Town prints:
      """
      Ran 17 shell commands.
      """
