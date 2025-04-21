Feature: display all executed Git commands

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | parent | feature | main   | local, origin |
      | child  | feature | parent | local, origin |
    And the current branch is "child"
    When I run "git-town set-parent main --verbose"

  Scenario: result
    And Git Town runs the commands
      | BRANCH | TYPE    | COMMAND                                          |
      |        | backend | git version                                      |
      |        | backend | git rev-parse --show-toplevel                    |
      |        | backend | git config -lz --includes --global               |
      |        | backend | git config -lz --includes --local                |
      |        | backend | git status -z --ignore-submodules                |
      |        | backend | git rev-parse -q --verify MERGE_HEAD             |
      |        | backend | git rev-parse -q --verify REBASE_HEAD            |
      |        | backend | git stash list                                   |
      |        | backend | git -c core.abbrev=40 branch -vva --sort=refname |
      |        | backend | git remote get-url origin                        |
      |        | backend | git config git-town-branch.child.parent main     |
      |        | backend | git -c core.abbrev=40 branch -vva --sort=refname |
      |        | backend | git config -lz --includes --global               |
      |        | backend | git config -lz --includes --local                |
      |        | backend | git stash list                                   |
    And Git Town prints:
      """
      Ran 15 shell commands.
      """
    And this lineage exists now
      | BRANCH | PARENT |
      | child  | main   |
      | parent | main   |

  Scenario: undo
    When I run "git-town undo --verbose"
    Then Git Town runs the commands
      | BRANCH | TYPE    | COMMAND                                          |
      |        | backend | git version                                      |
      |        | backend | git rev-parse --show-toplevel                    |
      |        | backend | git config -lz --includes --global               |
      |        | backend | git config -lz --includes --local                |
      |        | backend | git status -z --ignore-submodules                |
      |        | backend | git rev-parse -q --verify MERGE_HEAD             |
      |        | backend | git rev-parse -q --verify REBASE_HEAD            |
      |        | backend | git stash list                                   |
      |        | backend | git -c core.abbrev=40 branch -vva --sort=refname |
      |        | backend | git remote get-url origin                        |
      |        | backend | git rev-parse --verify --abbrev-ref @{-1}        |
      |        | backend | git remote get-url origin                        |
      |        | backend | git config git-town-branch.child.parent parent   |
    And Git Town prints:
      """
      Ran 13 shell commands.
      """
    And the initial commits exist now
    And the initial branches and lineage exist now
