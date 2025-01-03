Feature: display all executed Git commands with uncommitted changes

  Background:
    Given a Git repo with origin
    And the commits
      | BRANCH | LOCATION | MESSAGE     |
      | main   | origin   | main commit |
    And the current branch is "main"

  Scenario: result
    When I run "git-town hack new --verbose"
    Then Git Town runs the commands
      | BRANCH | TYPE     | COMMAND                                       |
      |        | backend  | git version                                   |
      |        | backend  | git rev-parse --show-toplevel                 |
      |        | backend  | git config -lz --includes --global            |
      |        | backend  | git config -lz --includes --local             |
      |        | backend  | git branch -vva --sort=refname                |
      |        | backend  | git rev-parse --verify --abbrev-ref @{-1}     |
      |        | backend  | git status --long --ignore-submodules         |
      |        | backend  | git stash list                                |
      |        | backend  | git branch -vva --sort=refname                |
      |        | backend  | git remote get-url origin                     |
      |        | backend  | git remote                                    |
      | main   | frontend | git add -A                                    |
      |        | frontend | git stash -m "Git Town WIP"                   |
      |        | backend  | git show-ref --verify --quiet refs/heads/main |
      | main   | frontend | git checkout -b new                           |
      |        | backend  | git show-ref --verify --quiet refs/heads/main |
      |        | backend  | git config git-town-branch.new.parent main    |
      |        | backend  | git show-ref --verify --quiet refs/heads/main |
      |        | backend  | git stash list                                |
      |        | backend  | git branch -vva --sort=refname                |
      |        | backend  | git config -lz --includes --global            |
      |        | backend  | git config -lz --includes --local             |
      |        | backend  | git stash list                                |
    And Git Town prints:
      """
      Ran 24 shell commands.
      """
    And the current branch is now "new"

  Scenario: undo
    Given I ran "git-town hack new"
    When I run "git-town undo --verbose"
    Then Git Town runs the commands
      | BRANCH | TYPE     | COMMAND                                       |
      |        | backend  | git version                                   |
      |        | backend  | git rev-parse --show-toplevel                 |
      |        | backend  | git config -lz --includes --global            |
      |        | backend  | git config -lz --includes --local             |
      |        | backend  | git status --long --ignore-submodules         |
      |        | backend  | git stash list                                |
      |        | backend  | git branch -vva --sort=refname                |
      |        | backend  | git remote get-url origin                     |
      |        | backend  | git rev-parse --verify --abbrev-ref @{-1}     |
      |        | backend  | git remote get-url origin                     |
      | new    | frontend | git add -A                                    |
      |        | frontend | git stash -m "Git Town WIP"                   |
      |        | frontend | git checkout main                             |
      | main   | frontend | git branch -D new                             |
      |        | backend  | git config --unset git-town-branch.new.parent |
      |        | backend  | git stash list                                |
    And Git Town prints:
      """
      Ran 17 shell commands.
      """
    And the current branch is now "main"
