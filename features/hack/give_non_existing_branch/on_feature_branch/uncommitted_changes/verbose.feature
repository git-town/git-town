Feature: display all executed Git commands with uncommitted changes

  Background:
    Given the commits
      | BRANCH | LOCATION | MESSAGE     |
      | main   | origin   | main commit |
    And the current branch is "main"
    And an uncommitted file

  Scenario: result
    When I run "git-town hack new --verbose"
    Then it runs the commands
      | BRANCH | TYPE     | COMMAND                                       |
      |        | backend  | git version                                   |
      |        | backend  | git config -lz --global                       |
      |        | backend  | git config -lz --local                        |
      |        | backend  | git rev-parse --show-toplevel                 |
      |        | backend  | git rev-parse --verify --abbrev-ref @{-1}     |
      |        | backend  | git status --long --ignore-submodules         |
      |        | backend  | git stash list                                |
      |        | backend  | git branch -vva --sort=refname                |
      |        | backend  | git remote                                    |
      | main   | frontend | git add -A                                    |
      |        | frontend | git stash                                     |
      |        | backend  | git show-ref --verify --quiet refs/heads/main |
      | main   | frontend | git checkout -b new                           |
      |        | backend  | git show-ref --verify --quiet refs/heads/main |
      |        | backend  | git config git-town-branch.new.parent main    |
      |        | backend  | git show-ref --verify --quiet refs/heads/main |
      |        | backend  | git stash list                                |
      | new    | frontend | git stash pop                                 |
      |        | backend  | git branch -vva --sort=refname                |
      |        | backend  | git config -lz --global                       |
      |        | backend  | git config -lz --local                        |
      |        | backend  | git stash list                                |
    And it prints:
      """
      Ran 22 shell commands.
      """
    And the current branch is now "new"
    And the uncommitted file still exists

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
      | new    | frontend | git add -A                                    |
      |        | frontend | git stash                                     |
      |        | frontend | git checkout main                             |
      | main   | frontend | git branch -D new                             |
      |        | backend  | git config --unset git-town-branch.new.parent |
      |        | backend  | git stash list                                |
      | main   | frontend | git stash pop                                 |
    And it prints:
      """
      Ran 16 shell commands.
      """
    And the current branch is now "main"
    And the uncommitted file still exists
