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
      | BRANCH | TYPE    | COMMAND                                                                                                                                                                                                                                                                                                                                          |
      |        | backend | git version                                                                                                                                                                                                                                                                                                                                      |
      |        | backend | git rev-parse --show-toplevel                                                                                                                                                                                                                                                                                                                    |
      |        | backend | git config -lz --global                                                                                                                                                                                                                                                                                                                          |
      |        | backend | git config -lz --local                                                                                                                                                                                                                                                                                                                           |
      |        | backend | git config -lz                                                                                                                                                                                                                                                                                                                                   |
      |        | backend | git status -z --ignore-submodules                                                                                                                                                                                                                                                                                                                |
      |        | backend | git rev-parse --verify -q MERGE_HEAD                                                                                                                                                                                                                                                                                                             |
      |        | backend | git rev-parse --absolute-git-dir                                                                                                                                                                                                                                                                                                                 |
      |        | backend | git remote get-url origin                                                                                                                                                                                                                                                                                                                        |
      |        | backend | git stash list                                                                                                                                                                                                                                                                                                                                   |
      |        | backend | git for-each-ref "--format=refname:%(refname) branchname:%(refname:lstrip=2) sha:%(objectname) head:%(if)%(HEAD)%(then)Y%(else)N%(end) worktree:%(if)%(worktreepath)%(then)Y%(else)N%(end) symref:%(if)%(symref)%(then)Y%(else)N%(end) upstream:%(upstream:lstrip=2) track:%(upstream:track,nobracket)" --sort=refname refs/heads/ refs/remotes/ |
      |        | backend | git remote                                                                                                                                                                                                                                                                                                                                       |
      |        | backend | git config git-town-branch.child.parent main                                                                                                                                                                                                                                                                                                     |
      |        | backend | git for-each-ref "--format=refname:%(refname) branchname:%(refname:lstrip=2) sha:%(objectname) head:%(if)%(HEAD)%(then)Y%(else)N%(end) worktree:%(if)%(worktreepath)%(then)Y%(else)N%(end) symref:%(if)%(symref)%(then)Y%(else)N%(end) upstream:%(upstream:lstrip=2) track:%(upstream:track,nobracket)" --sort=refname refs/heads/ refs/remotes/ |
      |        | backend | git config -lz --global                                                                                                                                                                                                                                                                                                                          |
      |        | backend | git config -lz --local                                                                                                                                                                                                                                                                                                                           |
      |        | backend | git stash list                                                                                                                                                                                                                                                                                                                                   |
    And Git Town prints:
      """
      Ran 17 shell commands.
      """
    And this lineage exists now
      """
      main
        child
        parent
      """

  Scenario: undo
    When I run "git-town undo --verbose"
    Then Git Town runs the commands
      | BRANCH | TYPE    | COMMAND                                                                                                                                                                                                                                                                                                                                          |
      |        | backend | git version                                                                                                                                                                                                                                                                                                                                      |
      |        | backend | git rev-parse --show-toplevel                                                                                                                                                                                                                                                                                                                    |
      |        | backend | git config -lz --global                                                                                                                                                                                                                                                                                                                          |
      |        | backend | git config -lz --local                                                                                                                                                                                                                                                                                                                           |
      |        | backend | git config -lz                                                                                                                                                                                                                                                                                                                                   |
      |        | backend | git status -z --ignore-submodules                                                                                                                                                                                                                                                                                                                |
      |        | backend | git rev-parse --verify -q MERGE_HEAD                                                                                                                                                                                                                                                                                                             |
      |        | backend | git rev-parse --absolute-git-dir                                                                                                                                                                                                                                                                                                                 |
      |        | backend | git remote get-url origin                                                                                                                                                                                                                                                                                                                        |
      |        | backend | git stash list                                                                                                                                                                                                                                                                                                                                   |
      |        | backend | git for-each-ref "--format=refname:%(refname) branchname:%(refname:lstrip=2) sha:%(objectname) head:%(if)%(HEAD)%(then)Y%(else)N%(end) worktree:%(if)%(worktreepath)%(then)Y%(else)N%(end) symref:%(if)%(symref)%(then)Y%(else)N%(end) upstream:%(upstream:lstrip=2) track:%(upstream:track,nobracket)" --sort=refname refs/heads/ refs/remotes/ |
      |        | backend | git remote                                                                                                                                                                                                                                                                                                                                       |
      |        | backend | git rev-parse --verify --abbrev-ref @{-1}                                                                                                                                                                                                                                                                                                        |
      |        | backend | git config git-town-branch.child.parent parent                                                                                                                                                                                                                                                                                                   |
    And Git Town prints:
      """
      Ran 14 shell commands.
      """
    And the initial branches and lineage exist now
    And the initial commits exist now
