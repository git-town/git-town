Feature: merging a branch verbosely

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE      | FILE NAME  | FILE CONTENT  |
      | alpha  | local, origin | alpha commit | alpha-file | alpha content |
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | beta | feature | alpha  | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT |
      | beta   | local, origin | beta commit | beta-file | beta content |
    And the current branch is "beta"
    And Git setting "git-town.sync-feature-strategy" is "merge"
    When I run "git-town merge -v"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                          |
      |        | git version                                      |
      |        | git rev-parse --show-toplevel                    |
      |        | git config -lz --includes --global               |
      |        | git config -lz --includes --local                |
      |        | git -c core.abbrev=40 branch -vva --sort=refname |
      |        | git status -z --ignore-submodules                |
      |        | git rev-parse --verify -q MERGE_HEAD             |
      |        | git rev-parse --absolute-git-dir                 |
      |        | git remote                                       |
      | beta   | git fetch --prune --tags                         |
      | (none) | git stash list                                   |
      |        | git -c core.abbrev=40 branch -vva --sort=refname |
      |        | git remote get-url origin                        |
      |        | git rev-parse --verify --abbrev-ref @{-1}        |
      |        | git log alpha..beta --format=%s --reverse        |
      |        | git log main..alpha --format=%s --reverse        |
      |        | git log --no-merges alpha ^beta                  |
      |        | git config git-town-branch.beta.parent main      |
      |        | git config --unset git-town-branch.alpha.parent  |
      | beta   | git branch -D alpha                              |
      |        | git push origin :alpha                           |
      | (none) | git -c core.abbrev=40 branch -vva --sort=refname |
      |        | git config -lz --includes --global               |
      |        | git config -lz --includes --local                |
      |        | git stash list                                   |
    And Git Town prints:
      """
      Ran 25 shell commands.
      """

  Scenario: undo
    When I run "git-town undo -v"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                          |
      |        | git version                                      |
      |        | git rev-parse --show-toplevel                    |
      |        | git config -lz --includes --global               |
      |        | git config -lz --includes --local                |
      |        | git status -z --ignore-submodules                |
      |        | git rev-parse --verify -q MERGE_HEAD             |
      |        | git rev-parse --absolute-git-dir                 |
      |        | git stash list                                   |
      |        | git -c core.abbrev=40 branch -vva --sort=refname |
      |        | git remote get-url origin                        |
      |        | git rev-parse --verify --abbrev-ref @{-1}        |
      |        | git remote get-url origin                        |
      | beta   | git branch alpha {{ sha 'alpha commit' }}        |
      |        | git push -u origin alpha                         |
      | (none) | git config git-town-branch.alpha.parent main     |
      |        | git config git-town-branch.beta.parent alpha     |
    And the initial commits exist now
    And the initial lineage exists now
