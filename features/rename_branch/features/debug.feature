Feature: display debug statistics

  Background:
    Given the current branch is a feature branch "old"
    And the commits
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
      | old    | local, origin | old commit  |

  Scenario: result
    When I run "git-town rename-branch new --debug"
    Then it runs the debug commands
      | git config -lz --local                        |
      | git config -lz --global                       |
      | git rev-parse                                 |
      | git rev-parse --show-toplevel                 |
      | git version                                   |
      | git branch -a                                 |
      | git status                                    |
      | git rev-parse --abbrev-ref HEAD               |
      | git branch                                    |
      | git branch -r                                 |
      | git rev-parse old                             |
      | git rev-parse origin/old                      |
      | git branch -a                                 |
      | git rev-parse --verify --abbrev-ref @{-1}     |
      | git status --porcelain --ignore-submodules    |
      | git config --unset git-town-branch.old.parent |
      | git config git-town-branch.new.parent main    |
      | git rev-parse old                             |
      | git log main..old                             |
      | git branch                                    |
      | git branch                                    |
      | git rev-parse --verify --abbrev-ref @{-1}     |
      | git checkout main                             |
      | git checkout new                              |
    And the current branch is now "new"

  Scenario: undo
    Given I ran "git-town rename-branch new"
    When I run "git-town undo --debug"
    Then it runs the debug commands
      | git config -lz --local                        |
      | git config -lz --global                       |
      | git rev-parse                                 |
      | git rev-parse --show-toplevel                 |
      | git version                                   |
      | git branch -a                                 |
      | git status                                    |
      | git rev-parse --abbrev-ref HEAD               |
      | git rev-parse origin/new                      |
      | git config --unset git-town-branch.new.parent |
      | git config git-town-branch.old.parent main    |
      | git rev-parse new                             |
      | git log old..new                              |
    And the current branch is now "old"
