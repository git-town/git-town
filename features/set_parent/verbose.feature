@messyoutput
Feature: display all executed Git commands

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | parent | feature | main   | local, origin |
      | child  | feature | parent | local, origin |
    And the current branch is "child"
    When I run "git-town set-parent --verbose" and enter into the dialog:
      | DIALOG                 | KEYS     |
      | parent branch of child | up enter |

  Scenario: result
    Then it prints:
      """
      Selected parent branch for "child": main
      """
    And it runs the commands
      | BRANCH | TYPE    | COMMAND                                      |
      |        | backend | git version                                  |
      |        | backend | git rev-parse --show-toplevel                |
      |        | backend | git config -lz --includes --global           |
      |        | backend | git config -lz --includes --local            |
      |        | backend | git status --long --ignore-submodules        |
      |        | backend | git stash list                               |
      |        | backend | git branch -vva --sort=refname               |
      |        | backend | git remote get-url origin                    |
      |        | backend | git config git-town-branch.child.parent main |
      |        | backend | git branch -vva --sort=refname               |
      |        | backend | git config -lz --includes --global           |
      |        | backend | git config -lz --includes --local            |
      |        | backend | git stash list                               |
    And it prints:
      """
      Ran 13 shell commands.
      """
    And this lineage exists now
      | BRANCH | PARENT |
      | child  | main   |
      | parent | main   |
    And the current branch is still "child"

  Scenario: undo
    When I run "git-town undo --verbose"
    Then it runs the commands
      | BRANCH | TYPE    | COMMAND                                        |
      |        | backend | git version                                    |
      |        | backend | git rev-parse --show-toplevel                  |
      |        | backend | git config -lz --includes --global             |
      |        | backend | git config -lz --includes --local              |
      |        | backend | git status --long --ignore-submodules          |
      |        | backend | git stash list                                 |
      |        | backend | git branch -vva --sort=refname                 |
      |        | backend | git remote get-url origin                      |
      |        | backend | git rev-parse --verify --abbrev-ref @{-1}      |
      |        | backend | git remote get-url origin                      |
      |        | backend | git config git-town-branch.child.parent parent |
    And it prints:
      """
      Ran 11 shell commands.
      """
    And the current branch is still "child"
    And the initial commits exist now
    And the initial branches and lineage exist now
