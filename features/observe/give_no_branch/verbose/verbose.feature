Feature: observe the current branch verbosely

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the current branch is "feature"
    When I run "git-town observe --verbose"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                |
      |        | git version                                            |
      |        | git rev-parse --show-toplevel                          |
      |        | git config -lz --includes --global                     |
      |        | git config -lz --includes --local                      |
      |        | git -c core.abbrev=40 branch -vva --sort=refname       |
      |        | git config git-town-branch.feature.branchtype observed |
      |        | git -c core.abbrev=40 branch -vva --sort=refname       |
      |        | git config -lz --includes --global                     |
      |        | git config -lz --includes --local                      |
    And Git Town prints:
      """
      Ran 9 shell commands
      """
    And Git Town prints:
      """
      branch "feature" is now an observed branch
      """
    And branch "feature" now has type "observed"

  Scenario: undo
    When I run "git-town undo --verbose"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                               |
      |        | git version                                           |
      |        | git rev-parse --show-toplevel                         |
      |        | git config -lz --includes --global                    |
      |        | git config -lz --includes --local                     |
      |        | git status -z --ignore-submodules                     |
      |        | git rev-parse -q --verify MERGE_HEAD                  |
      |        | git rev-parse --absolute-git-dir                      |
      |        | git stash list                                        |
      |        | git -c core.abbrev=40 branch -vva --sort=refname      |
      |        | git remote get-url origin                             |
      |        | git rev-parse --verify --abbrev-ref @{-1}             |
      |        | git remote get-url origin                             |
      |        | git config --unset git-town-branch.feature.branchtype |
    And Git Town prints:
      """
      Ran 13 shell commands
      """
    And branch "feature" now has type "feature"
