Feature: make the current branch a contribution branch verbosely

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the current branch is "feature"
    When I run "git-town contribute --verbose"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                    |
      |        | git version                                                |
      |        | git rev-parse --show-toplevel                              |
      |        | git config -lz --includes --global                         |
      |        | git config -lz --includes --local                          |
      |        | git branch -vva --sort=refname                             |
      |        | git config git-town-branch.feature.branchtype contribution |
      |        | git branch -vva --sort=refname                             |
      |        | git config -lz --includes --global                         |
      |        | git config -lz --includes --local                          |
    And Git Town prints:
      """
      Ran 9 shell commands
      """
    And Git Town prints:
      """
      branch "feature" is now a contribution branch
      """
    And the current branch is still "feature"
    And branch "feature" now has type "contribution"

  Scenario: undo
    When I run "git-town undo --verbose"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                               |
      |        | git version                                           |
      |        | git rev-parse --show-toplevel                         |
      |        | git config -lz --includes --global                    |
      |        | git config -lz --includes --local                     |
      |        | git status --long --ignore-submodules                 |
      |        | git stash list                                        |
      |        | git branch -vva --sort=refname                        |
      |        | git remote get-url origin                             |
      |        | git rev-parse --verify --abbrev-ref @{-1}             |
      |        | git remote get-url origin                             |
      |        | git config --unset git-town-branch.feature.branchtype |
    And Git Town prints:
      """
      Ran 11 shell commands
      """
    And the current branch is still "feature"
    And branch "feature" now has type "feature"