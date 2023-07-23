Feature: display debug statistics

  Background:
    Given the feature branches "active" and "old"
    And the commits
      | BRANCH | LOCATION      | MESSAGE       |
      | active | local, origin | active commit |
      | old    | local, origin | old commit    |
    And origin deletes the "old" branch
    And the current branch is "old"

  Scenario: result
    When I run "git-town prune-branches --debug"
    Then it runs the debug commands
      | git config -lz --local                        |
      | git config -lz --global                       |
      | git rev-parse                                 |
      | git rev-parse --show-toplevel                 |
      | git version                                   |
      | git branch -a                                 |
      | git remote                                    |
      | git status                                    |
      | git rev-parse --abbrev-ref HEAD               |
      | git branch -vv                                |
      | git rev-parse --verify --abbrev-ref @{-1}     |
      | git status --porcelain --ignore-submodules    |
      | git config --unset git-town-branch.old.parent |
      | git rev-parse old                             |
      | git log main..old                             |
      | git branch                                    |
      | git branch                                    |
      | git rev-parse --verify --abbrev-ref @{-1}     |
      | git checkout main                             |
      | git checkout main                             |
    And the current branch is now "main"
    And the branches are now
      | REPOSITORY    | BRANCHES     |
      | local, origin | main, active |
    And this branch lineage exists now
      | BRANCH | PARENT |
      | active | main   |

  Scenario: undo
    Given I ran "git-town prune-branches"
    When I run "git-town undo --debug"
    Then it runs the debug commands
      | git config -lz --local                     |
      | git config -lz --global                    |
      | git rev-parse                              |
      | git rev-parse --show-toplevel              |
      | git version                                |
      | git branch -a                              |
      | git status                                 |
      | git rev-parse --abbrev-ref HEAD            |
      | git config git-town-branch.old.parent main |
    And the current branch is now "old"
    And the initial branches and hierarchy exist
