Feature: display debug statistics

  Background:
    Given the current branch is a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |

  Scenario: result
    When I run "git-town ship -m done --debug"
    Then it runs the debug commands
      | git config -lz --local                            |
      | git config -lz --global                           |
      | git rev-parse                                     |
      | git rev-parse --show-toplevel                     |
      | git version                                       |
      | git branch -a                                     |
      | git remote get-url origin                         |
      | git remote get-url origin                         |
      | git remote get-url origin                         |
      | git remote get-url origin                         |
      | git remote                                        |
      | git status                                        |
      | git rev-parse --abbrev-ref HEAD                   |
      | git status --porcelain --ignore-submodules        |
      | git branch -r                                     |
      | git rev-parse --verify --abbrev-ref @{-1}         |
      | git status --porcelain --ignore-submodules        |
      | git rev-parse HEAD                                |
      | git rev-list --left-right main...origin/main      |
      | git rev-parse HEAD                                |
      | git rev-parse HEAD                                |
      | git diff main..feature                            |
      | git shortlog -s -n -e main..feature               |
      | git config user.name                              |
      | git config user.email                             |
      | git rev-parse HEAD                                |
      | git rev-list --left-right main...origin/main      |
      | git rev-parse feature                             |
      | git log main..feature                             |
      | git config --unset git-town-branch.feature.parent |
      | git branch                                        |
      | git rev-parse --verify --abbrev-ref @{-1}         |
      | git checkout main                                 |
      | git checkout main                                 |
    And the current branch is now "main"

  Scenario: undo
    Given I ran "git-town ship -m done"
    When I run "git-town undo --debug"
    Then it runs the debug commands
      | git config -lz --local                         |
      | git config -lz --global                        |
      | git rev-parse                                  |
      | git rev-parse --show-toplevel                  |
      | git version                                    |
      | git branch -a                                  |
      | git config git-town-branch.feature.parent main |
      | git status                                     |
      | git rev-parse --abbrev-ref HEAD                |
      | git rev-list --left-right main...origin/main   |
      | git rev-parse HEAD                             |
      | git rev-parse HEAD                             |
    And the current branch is now "feature"
