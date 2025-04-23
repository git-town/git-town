Feature: display all executed Git commands when using the always-merge strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
    And the current branch is "feature"
    And Git setting "git-town.ship-strategy" is "always-merge"
    When I run "git-town ship --verbose" and close the editor

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | TYPE     | COMMAND                                           |
      |         | backend  | git version                                       |
      |         | backend  | git rev-parse --show-toplevel                     |
      |         | backend  | git config -lz --includes --global                |
      |         | backend  | git config -lz --includes --local                 |
      |         | backend  | git status -z --ignore-submodules                 |
      |         | backend  | git rev-parse -q --verify MERGE_HEAD              |
      |         | backend  | git rev-parse --absolute-git-dir                  |
      |         | backend  | git remote                                        |
      |         | backend  | git branch --show-current                         |
      | feature | frontend | git fetch --prune --tags                          |
      |         | backend  | git stash list                                    |
      |         | backend  | git -c core.abbrev=40 branch -vva --sort=refname  |
      |         | backend  | git rev-parse --verify --abbrev-ref @{-1}         |
      |         | backend  | git remote get-url origin                         |
      |         | backend  | git shortlog -s -n -e main..feature               |
      |         | backend  | git diff --shortstat main feature                 |
      | feature | frontend | git checkout main                                 |
      | main    | frontend | git merge --no-ff --edit -- feature               |
      |         | backend  | git rev-parse --verify -q refs/heads/main         |
      |         | backend  | git rev-list --left-right main...origin/main      |
      | main    | frontend | git push                                          |
      |         | backend  | git config --unset git-town-branch.feature.parent |
      | main    | frontend | git push origin :feature                          |
      |         | frontend | git branch -D feature                             |
      |         | backend  | git rev-parse --verify -q refs/heads/feature      |
      |         | backend  | git -c core.abbrev=40 branch -vva --sort=refname  |
      |         | backend  | git config -lz --includes --global                |
      |         | backend  | git config -lz --includes --local                 |
      |         | backend  | git stash list                                    |
    And Git Town prints:
      """
      Ran 29 shell commands.
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
      | main   | frontend | git branch feature {{ sha 'feature commit' }}    |
      |        | frontend | git push -u origin feature                       |
      |        | backend  | git rev-parse --verify -q refs/heads/feature     |
      | main   | frontend | git checkout feature                             |
      |        | backend  | git config git-town-branch.feature.parent main   |
    And Git Town prints:
      """
      Ran 17 shell commands.
      """
