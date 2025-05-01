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
      | BRANCH | TYPE    | COMMAND                                                                                                                                                                                                                                                                                                                                        |
      |        | backend | git version                                                                                                                                                                                                                                                                                                                                    |
      |        | backend | git rev-parse --show-toplevel                                                                                                                                                                                                                                                                                                                  |
      |        | backend | git config -lz --includes --global                                                                                                                                                                                                                                                                                                             |
      |        | backend | git config -lz --includes --local                                                                                                                                                                                                                                                                                                              |
      |        | backend | git config -lz --includes                                                                                                                                                                                                                                                                                                                      |
      |        | backend | git status -z --ignore-submodules                                                                                                                                                                                                                                                                                                              |
      |        | backend | git rev-parse --verify -q MERGE_HEAD                                                                                                                                                                                                                                                                                                           |
      |        | backend | git rev-parse --absolute-git-dir                                                                                                                                                                                                                                                                                                               |
      |        | backend | git stash list                                                                                                                                                                                                                                                                                                                                 |
      |        | backend | git for-each-ref --format=refname:%(refname) branchname:%(refname:lstrip=2) sha:%(objectname) head:%(if)%(HEAD)%(then)Y%(else)N%(end) worktree:%(if)%(worktreepath)%(then)Y%(else)N%(end) symref:%(if)%(symref)%(then)Y%(else)N%(end) upstream:%(upstream:lstrip=2) track:%(upstream:track,nobracket) --sort=refname refs/heads/ refs/remotes/ |
      |        | backend | git remote get-url origin                                                                                                                                                                                                                                                                                                                      |
      |        | backend | git config git-town-branch.child.parent main                                                                                                                                                                                                                                                                                                   |
      |        | backend | git for-each-ref --format=refname:%(refname) branchname:%(refname:lstrip=2) sha:%(objectname) head:%(if)%(HEAD)%(then)Y%(else)N%(end) worktree:%(if)%(worktreepath)%(then)Y%(else)N%(end) symref:%(if)%(symref)%(then)Y%(else)N%(end) upstream:%(upstream:lstrip=2) track:%(upstream:track,nobracket) --sort=refname refs/heads/ refs/remotes/ |
      |        | backend | git config -lz --includes --global                                                                                                                                                                                                                                                                                                             |
      |        | backend | git config -lz --includes --local                                                                                                                                                                                                                                                                                                              |
      |        | backend | git stash list                                                                                                                                                                                                                                                                                                                                 |
    And Git Town prints:
      """
      Ran 16 shell commands.
      """
    And this lineage exists now
      | BRANCH | PARENT |
      | child  | main   |
      | parent | main   |

  Scenario: undo
    When I run "git-town undo --verbose"
    Then Git Town runs the commands
      | BRANCH | TYPE    | COMMAND                                                                                                                                                                                                                                                                                                                                        |
      |        | backend | git version                                                                                                                                                                                                                                                                                                                                    |
      |        | backend | git rev-parse --show-toplevel                                                                                                                                                                                                                                                                                                                  |
      |        | backend | git config -lz --includes --global                                                                                                                                                                                                                                                                                                             |
      |        | backend | git config -lz --includes --local                                                                                                                                                                                                                                                                                                              |
      |        | backend | git config -lz --includes                                                                                                                                                                                                                                                                                                                      |
      |        | backend | git status -z --ignore-submodules                                                                                                                                                                                                                                                                                                              |
      |        | backend | git rev-parse --verify -q MERGE_HEAD                                                                                                                                                                                                                                                                                                           |
      |        | backend | git rev-parse --absolute-git-dir                                                                                                                                                                                                                                                                                                               |
      |        | backend | git stash list                                                                                                                                                                                                                                                                                                                                 |
      |        | backend | git for-each-ref --format=refname:%(refname) branchname:%(refname:lstrip=2) sha:%(objectname) head:%(if)%(HEAD)%(then)Y%(else)N%(end) worktree:%(if)%(worktreepath)%(then)Y%(else)N%(end) symref:%(if)%(symref)%(then)Y%(else)N%(end) upstream:%(upstream:lstrip=2) track:%(upstream:track,nobracket) --sort=refname refs/heads/ refs/remotes/ |
      |        | backend | git remote get-url origin                                                                                                                                                                                                                                                                                                                      |
      |        | backend | git rev-parse --verify --abbrev-ref @{-1}                                                                                                                                                                                                                                                                                                      |
      |        | backend | git remote get-url origin                                                                                                                                                                                                                                                                                                                      |
      |        | backend | git config git-town-branch.child.parent parent                                                                                                                                                                                                                                                                                                 |
    And Git Town prints:
      """
      Ran 14 shell commands.
      """
    And the initial commits exist now
    And the initial branches and lineage exist now
