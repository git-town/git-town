Feature: observe a branch verbosely

  Background:
    Given the current branch is a feature branch "branch"
    And an uncommitted file
    When I run "git-town observe --verbose"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                                      |
      |        | git version                                  |
      |        | git config -lz --global                      |
      |        | git config -lz --local                       |
      |        | git rev-parse --show-toplevel                |
      |        | git status --long --ignore-submodules        |
      |        | git stash list                               |
      |        | git remote                                   |
      |        | git fetch --prune --tags                     |
      |        | git branch -vva --sort=refname               |
      |        | git config git-town.observed-branches branch |
      |        | git config -lz --global                      |
      |        | git config -lz --local                       |
    And it prints:
      """
      Ran 12 shell commands
      """
    And it prints:
      """
      branch "branch" is now an observed branch
      """
    And the current branch is still "branch"
    And branch "branch" is now observed
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo --verbose"
    Then it runs the commands
      | BRANCH | COMMAND                                       |
      |        | git version                                   |
      |        | git config -lz --global                       |
      |        | git config -lz --local                        |
      |        | git rev-parse --show-toplevel                 |
      |        | git status --long --ignore-submodules         |
      |        | git stash list                                |
      |        | git branch -vva --sort=refname                |
      |        | git rev-parse --verify --abbrev-ref @{-1}     |
      |        | git remote get-url origin                     |
      | branch | git add -A                                    |
      |        | git stash                                     |
      | <none> | git config --unset git-town.observed-branches |
      |        | git show-ref --verify --quiet refs/heads/     |
      |        | git stash list                                |
      | branch | git stash pop                                 |
    And it prints:
      """
      Ran 15 shell commands
      """
    And the current branch is still "branch"
    And branch "branch" is now a feature branch
    And the uncommitted file still exists
