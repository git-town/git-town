Feature: merging a branch in a stack with its parent

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
      | beta  | feature | alpha  | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE      | FILE NAME  | FILE CONTENT  |
      | alpha  | local, origin | alpha commit | alpha-file | alpha content |
      | beta   | local, origin | beta commit  | beta-file  | beta content  |
    And the current branch is "beta"
    And Git Town setting "sync-feature-strategy" is "merge"
    When I run "git-town merge -v"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                         |
      |        | git version                                     |
      |        | git rev-parse --show-toplevel                   |
      |        | git config -lz --includes --global              |
      |        | git config -lz --includes --local               |
      |        | git branch -vva --sort=refname                  |
      |        | git status --long --ignore-submodules           |
      |        | git remote                                      |
      | beta   | git fetch --prune --tags                        |
      | <none> | git stash list                                  |
      |        | git branch -vva --sort=refname                  |
      |        | git remote get-url origin                       |
      |        | git rev-parse --verify --abbrev-ref @{-1}       |
      |        | git log main..alpha --format=%s --reverse       |
      |        | git log alpha..beta --format=%s --reverse       |
      | beta   | git checkout alpha                              |
      | alpha  | git merge --no-edit --ff origin/alpha           |
      |        | git checkout beta                               |
      | beta   | git merge --no-edit --ff alpha                  |
      |        | git merge --no-edit --ff origin/beta            |
      | <none> | git rev-list --left-right beta...origin/beta    |
      | beta   | git push                                        |
      | <none> | git config git-town-branch.beta.parent main     |
      |        | git config --unset git-town-branch.alpha.parent |
      | beta   | git branch -D alpha                             |
      |        | git push origin :alpha                          |
      | <none> | git show-ref --verify --quiet refs/heads/main   |
      |        | git checkout main                               |
      |        | git checkout beta                               |
      |        | git branch -vva --sort=refname                  |
      |        | git config -lz --includes --global              |
      |        | git config -lz --includes --local               |
      |        | git stash list                                  |
    And Git Town prints:
      """
      Ran 32 shell commands.
      """
    And the current branch is still "beta"
    And this lineage exists now
      | BRANCH | PARENT |
      | beta   | main   |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                        | FILE NAME  | FILE CONTENT  |
      | beta   | local, origin | beta commit                    | beta-file  | beta content  |
      |        |               | alpha commit                   | alpha-file | alpha content |
      |        |               | Merge branch 'alpha' into beta |            |               |
    And these committed files exist now
      | BRANCH | NAME       | CONTENT       |
      | beta   | alpha-file | alpha content |
      |        | beta-file  | beta content  |

  Scenario: undo
    When I run "git-town undo -v"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                         |
      |        | git version                                     |
      |        | git rev-parse --show-toplevel                   |
      |        | git config -lz --includes --global              |
      |        | git config -lz --includes --local               |
      |        | git status --long --ignore-submodules           |
      |        | git stash list                                  |
      |        | git branch -vva --sort=refname                  |
      |        | git remote get-url origin                       |
      |        | git rev-parse --verify --abbrev-ref @{-1}       |
      |        | git remote get-url origin                       |
      |        | git rev-parse --short HEAD                      |
      | beta   | git reset --hard {{ sha 'beta commit' }}        |
      | <none> | git rev-list --left-right beta...origin/beta    |
      | beta   | git push --force-with-lease --force-if-includes |
      |        | git branch alpha {{ sha 'alpha commit' }}       |
      |        | git push -u origin alpha                        |
      | <none> | git config git-town-branch.alpha.parent main    |
      |        | git config git-town-branch.beta.parent alpha    |
    And the current branch is still "beta"
    And the initial commits exist now
    And the initial lineage exists now
