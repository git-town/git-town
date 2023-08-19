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
    Then it runs the commands
      | BRANCH | TYPE     | COMMAND                                       |
      |        | backend  | git version                                   |
      |        | backend  | git config -lz --local                        |
      |        | backend  | git config -lz --global                       |
      |        | backend  | git rev-parse --show-toplevel                 |
      |        | backend  | git remote                                    |
      |        | backend  | git status                                    |
      |        | backend  | git rev-parse --abbrev-ref HEAD               |
      | old    | frontend | git fetch --prune --tags                      |
      |        | backend  | git branch -vva                               |
      |        | backend  | git rev-parse --verify --abbrev-ref @{-1}     |
      | old    | frontend | git checkout main                             |
      |        | backend  | git config --unset git-town-branch.old.parent |
      |        | backend  | git rev-parse old                             |
      |        | backend  | git log main..old                             |
      | main   | frontend | git branch -D old                             |
      |        | backend  | git show-ref --quiet refs/heads/main          |
      |        | backend  | git show-ref --quiet refs/heads/old           |
      |        | backend  | git rev-parse --verify --abbrev-ref @{-1}     |
      |        | backend  | git checkout main                             |
      |        | backend  | git checkout main                             |
    And it prints:
      """
      Ran 20 shell commands.
      """
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
    Then it runs the commands
      | BRANCH | TYPE     | COMMAND                                    |
      |        | backend  | git version                                |
      |        | backend  | git config -lz --local                     |
      |        | backend  | git config -lz --global                    |
      |        | backend  | git rev-parse --show-toplevel              |
      |        | backend  | git branch -vva                            |
      | main   | frontend | git branch old {{ sha 'old commit' }}      |
      |        | backend  | git config git-town-branch.old.parent main |
      | main   | frontend | git checkout old                           |
    And it prints:
      """
      Ran 8 shell commands.
      """
    And the current branch is now "old"
    And the initial branches and hierarchy exist
