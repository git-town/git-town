Feature: display debug statistics

  Background:
    Given the current branch is a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE               |
      | main    | local    | local main commit     |
      |         | origin   | origin main commit    |
      | feature | local    | local feature commit  |
      |         | origin   | origin feature commit |

  Scenario: result
    When I run "git-town sync --debug"
    Then it runs the commands
      | BRANCH  | TYPE     | COMMAND                                            |
      |         | backend  | git version                                        |
      |         | backend  | git config -lz --local                             |
      |         | backend  | git config -lz --global                            |
      |         | backend  | git rev-parse --show-toplevel                      |
      |         | backend  | git branch -vva                                    |
      |         | backend  | git remote                                         |
      | feature | frontend | git fetch --prune --tags                           |
      |         | backend  | git branch -vva                                    |
      |         | backend  | git rev-parse --verify --abbrev-ref @{-1}          |
      |         | backend  | git status --porcelain --ignore-submodules         |
      | feature | frontend | git checkout main                                  |
      |         | backend  | git rev-parse HEAD                                 |
      | main    | frontend | git rebase origin/main                             |
      |         | backend  | git rev-list --left-right main...origin/main       |
      | main    | frontend | git push                                           |
      |         | frontend | git checkout feature                               |
      |         | backend  | git rev-parse HEAD                                 |
      | feature | frontend | git merge --no-edit origin/feature                 |
      |         | backend  | git rev-parse HEAD                                 |
      | feature | frontend | git merge --no-edit main                           |
      |         | backend  | git rev-list --left-right feature...origin/feature |
      | feature | frontend | git push                                           |
      |         | backend  | git show-ref --quiet refs/heads/main               |
      |         | backend  | git rev-parse --verify --abbrev-ref @{-1}          |
    And it prints:
      """
      Ran 24 shell commands.
      """
    And all branches are now synchronized
