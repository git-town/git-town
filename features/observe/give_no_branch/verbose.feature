Feature: observe the current branch verbosely

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the current branch is "feature"
    And an uncommitted file
    When I run "git-town observe --verbose"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                                       |
      |        | git version                                   |
      |        | git rev-parse --show-toplevel                 |
      |        | git config -lz --includes --global            |
      |        | git config -lz --includes --local             |
      |        | git branch -vva --sort=refname                |
      |        | git config git-town.observed-branches feature |
      |        | git branch -vva --sort=refname                |
      |        | git config -lz --includes --global            |
      |        | git config -lz --includes --local             |
    And it prints:
      """
      Ran 9 shell commands
      """
    And it prints:
      """
      branch "feature" is now an observed branch
      """
    And the current branch is still "feature"
    And branch "feature" is now observed
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo --verbose"
    Then it runs the commands
      | BRANCH  | COMMAND                                       |
      |         | git version                                   |
      |         | git rev-parse --show-toplevel                 |
      |         | git config -lz --includes --global            |
      |         | git config -lz --includes --local             |
      |         | git status --long --ignore-submodules         |
      |         | git stash list                                |
      |         | git branch -vva --sort=refname                |
      |         | git remote get-url origin                     |
      |         | git rev-parse --verify --abbrev-ref @{-1}     |
      | feature | git add -A                                    |
      |         | git stash                                     |
      | <none>  | git config --unset git-town.observed-branches |
      |         | git stash list                                |
      | feature | git stash pop                                 |
    And it prints:
      """
      Ran 14 shell commands
      """
    And the current branch is still "feature"
    And branch "feature" is now a feature branch
    And the uncommitted file still exists
