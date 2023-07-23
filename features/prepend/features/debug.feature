Feature: display debug statistics

  Background:
    Given the current branch is a feature branch "old"
    And the commits
      | BRANCH | LOCATION      | MESSAGE    |
      | old    | local, origin | old commit |

  Scenario: result
    When I run "git-town prepend parent --debug"
    Then it runs the debug commands
      | git config -lz --local                        |
      | git config -lz --global                       |
      | git rev-parse                                 |
      | git rev-parse --show-toplevel                 |
      | git version                                   |
      | git branch -a                                 |
      | git status                                    |
      | git rev-parse --abbrev-ref HEAD               |
      | git remote                                    |
      | git branch -a                                 |
      | git branch -r                                 |
      | git rev-parse --verify --abbrev-ref @{-1}     |
      | git status --porcelain --ignore-submodules    |
      | git rev-parse HEAD                            |
      | git rev-list --left-right main...origin/main  |
      | git config git-town-branch.parent.parent main |
      | git config git-town-branch.old.parent parent  |
      | git branch                                    |
      | git branch                                    |
      | git rev-parse --verify --abbrev-ref @{-1}     |
      | git checkout old                              |
      | git checkout parent                           |
    And the current branch is now "parent"

  Scenario: undo
    Given I ran "git-town prepend parent"
    When I run "git-town undo --debug"
    Then it runs the debug commands
      | git config -lz --local                           |
      | git config -lz --global                          |
      | git rev-parse                                    |
      | git rev-parse --show-toplevel                    |
      | git version                                      |
      | git branch -a                                    |
      | git status                                       |
      | git rev-parse --abbrev-ref HEAD                  |
      | git config git-town-branch.old.parent main       |
      | git config --unset git-town-branch.parent.parent |
      | git rev-parse parent                             |
      | git log main..parent                             |
    And the current branch is now "old"
