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
      | BRANCH | COMMAND                                                                                                                                                                                                                                                                                                                                          |
      |        | git version                                                                                                                                                                                                                                                                                                                                      |
      |        | git rev-parse --show-toplevel                                                                                                                                                                                                                                                                                                                    |
      |        | git config -lz --global                                                                                                                                                                                                                                                                                                                          |
      |        | git config -lz --local                                                                                                                                                                                                                                                                                                                           |
      |        | git config -lz                                                                                                                                                                                                                                                                                                                                   |
      |        | git for-each-ref "--format=refname:%(refname) branchname:%(refname:lstrip=2) sha:%(objectname) head:%(if)%(HEAD)%(then)Y%(else)N%(end) worktree:%(if)%(worktreepath)%(then)Y%(else)N%(end) symref:%(if)%(symref)%(then)Y%(else)N%(end) upstream:%(upstream:lstrip=2) track:%(upstream:track,nobracket)" --sort=refname refs/heads/ refs/remotes/ |
      |        | git config git-town-branch.feature.branchtype contribution                                                                                                                                                                                                                                                                                       |
      |        | git for-each-ref "--format=refname:%(refname) branchname:%(refname:lstrip=2) sha:%(objectname) head:%(if)%(HEAD)%(then)Y%(else)N%(end) worktree:%(if)%(worktreepath)%(then)Y%(else)N%(end) symref:%(if)%(symref)%(then)Y%(else)N%(end) upstream:%(upstream:lstrip=2) track:%(upstream:track,nobracket)" --sort=refname refs/heads/ refs/remotes/ |
      |        | git config -lz --global                                                                                                                                                                                                                                                                                                                          |
      |        | git config -lz --local                                                                                                                                                                                                                                                                                                                           |
    And Git Town prints:
      """
      Ran 10 shell commands
      """
    And Git Town prints:
      """
      branch "feature" is now a contribution branch
      """
    And branch "feature" now has type "contribution"

  Scenario: undo
    When I run "git-town undo --verbose"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                                                                                                                                                                                                                                                                                                          |
      |        | git version                                                                                                                                                                                                                                                                                                                                      |
      |        | git rev-parse --show-toplevel                                                                                                                                                                                                                                                                                                                    |
      |        | git config -lz --global                                                                                                                                                                                                                                                                                                                          |
      |        | git config -lz --local                                                                                                                                                                                                                                                                                                                           |
      |        | git config -lz                                                                                                                                                                                                                                                                                                                                   |
      |        | git status -z --ignore-submodules                                                                                                                                                                                                                                                                                                                |
      |        | git rev-parse --verify -q MERGE_HEAD                                                                                                                                                                                                                                                                                                             |
      |        | git rev-parse --absolute-git-dir                                                                                                                                                                                                                                                                                                                 |
      |        | git stash list                                                                                                                                                                                                                                                                                                                                   |
      |        | git for-each-ref "--format=refname:%(refname) branchname:%(refname:lstrip=2) sha:%(objectname) head:%(if)%(HEAD)%(then)Y%(else)N%(end) worktree:%(if)%(worktreepath)%(then)Y%(else)N%(end) symref:%(if)%(symref)%(then)Y%(else)N%(end) upstream:%(upstream:lstrip=2) track:%(upstream:track,nobracket)" --sort=refname refs/heads/ refs/remotes/ |
      |        | git remote get-url origin                                                                                                                                                                                                                                                                                                                        |
      |        | git rev-parse --verify --abbrev-ref @{-1}                                                                                                                                                                                                                                                                                                        |
      |        | git remote get-url origin                                                                                                                                                                                                                                                                                                                        |
      |        | git config --unset git-town-branch.feature.branchtype                                                                                                                                                                                                                                                                                            |
    And Git Town prints:
      """
      Ran 14 shell commands
      """
    And branch "feature" now has type "feature"
