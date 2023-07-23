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
    Then it runs the debug commands
      | git config -lz --local                             |
      | git config -lz --global                            |
      | git rev-parse                                      |
      | git rev-parse --show-toplevel                      |
      | git version                                        |
      | git branch -a                                      |
      | git remote                                         |
      | git status                                         |
      | git rev-parse --abbrev-ref HEAD                    |
      | git branch -r                                      |
      | git rev-parse --verify --abbrev-ref @{-1}          |
      | git status --porcelain --ignore-submodules         |
      | git rev-parse HEAD                                 |
      | git rev-list --left-right main...origin/main       |
      | git rev-parse HEAD                                 |
      | git rev-parse HEAD                                 |
      | git rev-list --left-right feature...origin/feature |
      | git branch                                         |
      | git branch                                         |
      | git rev-parse --verify --abbrev-ref @{-1}          |
    And all branches are now synchronized
