Feature: make a branch contribution verbosely

  Background:
    Given a Git repo clone
    And the branch
      | NAME   | TYPE    | PARENT | LOCATIONS |
      | branch | feature | main   | local     |
    And the current branch is "branch"
    And an uncommitted file
    When I run "git-town contribute --verbose"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                                          |
      |        | git version                                      |
      |        | git rev-parse --show-toplevel                    |
      |        | git config -lz --includes --global               |
      |        | git config -lz --includes --local                |
      |        | git branch -vva --sort=refname                   |
      |        | git config git-town.contribution-branches branch |
      |        | git branch -vva --sort=refname                   |
      |        | git config -lz --includes --global               |
      |        | git config -lz --includes --local                |
    And it prints:
      """
      Ran 9 shell commands
      """
    And it prints:
      """
      branch "branch" is now a contribution branch
      """
    And the current branch is still "branch"
    And branch "branch" is now a contribution branch
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo --verbose"
    Then it runs the commands
      | BRANCH | COMMAND                                           |
      |        | git version                                       |
      |        | git rev-parse --show-toplevel                     |
      |        | git config -lz --includes --global                |
      |        | git config -lz --includes --local                 |
      |        | git status --long --ignore-submodules             |
      |        | git stash list                                    |
      |        | git branch -vva --sort=refname                    |
      |        | git rev-parse --verify --abbrev-ref @{-1}         |
      |        | git remote get-url origin                         |
      | branch | git add -A                                        |
      |        | git stash                                         |
      | <none> | git config --unset git-town.contribution-branches |
      |        | git stash list                                    |
      | branch | git stash pop                                     |
    And it prints:
      """
      Ran 14 shell commands
      """
    And the current branch is still "branch"
    And branch "branch" is now a feature branch
    And the uncommitted file still exists
