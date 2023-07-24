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
      |         | backend  | git config -lz --local                             |
      |         | backend  | git config -lz --global                            |
      |         | backend  | git rev-parse                                      |
      |         | backend  | git rev-parse --show-toplevel                      |
      |         | backend  | git version                                        |
      |         | backend  | git branch -a                                      |
      |         | backend  | git remote                                         |
      |         | backend  | git status                                         |
      |         | backend  | git rev-parse --abbrev-ref HEAD                    |
      | feature | frontend | git fetch --prune --tags                           |
      |         | backend  | git branch -r                                      |
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
      |         | backend  | git branch                                         |
      |         | backend  | git branch                                         |
      |         | backend  | git rev-parse --verify --abbrev-ref @{-1}          |
    And all branches are now synchronized
