Feature: display debug statistics

  Background:
    Given the current branch is a feature branch "old"
    And the commits
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
      | old    | local, origin | old commit  |

  Scenario: result
    When I run "git-town rename-branch new --debug"
    Then it runs the commands
      | BRANCH | TYPE     | COMMAND                                       |
      |        | backend  | git version                                   |
      |        | backend  | git config -lz --local                        |
      |        | backend  | git config -lz --global                       |
      |        | backend  | git rev-parse --show-toplevel                 |
      |        | backend  | git branch -vva                               |
      |        | backend  | git remote                                    |
      | old    | frontend | git fetch --prune --tags                      |
      |        | backend  | git branch -vva                               |
      |        | backend  | git rev-parse --verify --abbrev-ref @{-1}     |
      | old    | frontend | git branch new old                            |
      |        | frontend | git checkout new                              |
      |        | backend  | git config --unset git-town-branch.old.parent |
      |        | backend  | git config git-town-branch.new.parent main    |
      | new    | frontend | git push -u origin new                        |
      |        | frontend | git push origin :old                          |
      |        | backend  | git rev-parse old                             |
      |        | backend  | git log main..old                             |
      | new    | frontend | git branch -D old                             |
      |        | backend  | git show-ref --quiet refs/heads/main          |
      |        | backend  | git show-ref --quiet refs/heads/old           |
      |        | backend  | git rev-parse --verify --abbrev-ref @{-1}     |
      |        | backend  | git checkout main                             |
      |        | backend  | git checkout new                              |
    And it prints:
      """
      Ran 23 shell commands.
      """
    And the current branch is now "new"

  Scenario: undo
    Given I ran "git-town rename-branch new"
    When I run "git-town undo --debug"
    Then it runs the commands
      | BRANCH | TYPE     | COMMAND                                       |
      |        | backend  | git version                                   |
      |        | backend  | git config -lz --local                        |
      |        | backend  | git config -lz --global                       |
      |        | backend  | git rev-parse --show-toplevel                 |
      |        | backend  | git branch -vva                               |
      | new    | frontend | git branch old {{ sha 'old commit' }}         |
      |        | frontend | git push -u origin old                        |
      |        | backend  | git rev-parse origin/new                      |
      | new    | frontend | git push origin :new                          |
      |        | backend  | git config --unset git-town-branch.new.parent |
      |        | backend  | git config git-town-branch.old.parent main    |
      | new    | frontend | git checkout old                              |
      |        | backend  | git rev-parse new                             |
      |        | backend  | git log old..new                              |
      | old    | frontend | git branch -D new                             |
    And it prints:
      """
      Ran 15 shell commands.
      """
    And the current branch is now "old"
