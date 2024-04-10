Feature: display all executed Git commands

  Background:
    Given the commits
      | BRANCH | LOCATION | MESSAGE     |
      | main   | origin   | main commit |
    And the current branch is "main"

  Scenario: result
    When I run "git-town hack new --verbose"
    Then it runs the commands
      | BRANCH | TYPE     | COMMAND                                       |
      |        | backend  | git version                                   |
      |        | backend  | git config -lz --global                       |
      |        | backend  | git config -lz --local                        |
      |        | backend  | git rev-parse --show-toplevel                 |
      |        | backend  | git status --long --ignore-submodules         |
      |        | backend  | git stash list                                |
      |        | backend  | git branch -vva --sort=refname                |
      |        | backend  | git remote                                    |
      | main   | frontend | git fetch --prune --tags                      |
      |        | backend  | git branch -vva --sort=refname                |
      |        | backend  | git rev-parse --verify --abbrev-ref @{-1}     |
      | main   | frontend | git rebase origin/main                        |
      |        | backend  | git rev-list --left-right main...origin/main  |
      |        | backend  | git show-ref --verify --quiet refs/heads/main |
      | main   | frontend | git branch new main                           |
      |        | backend  | git show-ref --verify --quiet refs/heads/main |
      |        | backend  | git config git-town-branch.new.parent main    |
      | main   | frontend | git checkout new                              |
      |        | backend  | git show-ref --verify --quiet refs/heads/main |
      |        | backend  | git branch -vva --sort=refname                |
      |        | backend  | git config -lz --global                       |
      |        | backend  | git config -lz --local                        |
      |        | backend  | git stash list                                |
    And it prints:
      """
      Ran 23 shell commands.
      """
    And the current branch is now "new"

  Scenario: undo
    Given I ran "git-town hack new"
    When I run "git-town undo --verbose"
    Then it runs the commands
      | BRANCH | TYPE     | COMMAND                                       |
      |        | backend  | git version                                   |
      |        | backend  | git config -lz --global                       |
      |        | backend  | git config -lz --local                        |
      |        | backend  | git rev-parse --show-toplevel                 |
      |        | backend  | git status --long --ignore-submodules         |
      |        | backend  | git stash list                                |
      |        | backend  | git branch -vva --sort=refname                |
      |        | backend  | git rev-parse --verify --abbrev-ref @{-1}     |
      |        | backend  | git remote get-url origin                     |
      | new    | frontend | git checkout main                             |
      |        | backend  | git rev-parse --short HEAD                    |
      | main   | frontend | git reset --hard {{ sha 'initial commit' }}   |
      |        | frontend | git branch -D new                             |
      |        | backend  | git config --unset git-town-branch.new.parent |
    And it prints:
      """
      Ran 14 shell commands.
      """
    And the current branch is now "main"
