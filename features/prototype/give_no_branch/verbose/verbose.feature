Feature: prototype the current branch verbosely

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    And the current branch is "feature"
    When I run "git-town prototype --verbose"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                                                                                                                                                                                                                                                                                                                                 |
      |        | git version                                                                                                                                                                                                                                                                                                                                                             |
      |        | git rev-parse --show-toplevel                                                                                                                                                                                                                                                                                                                                           |
      |        | git config -lz --includes --global                                                                                                                                                                                                                                                                                                                                      |
      |        | git config -lz --includes --local                                                                                                                                                                                                                                                                                                                                       |
      |        | git for-each-ref --format=refname:%(refname) branchname:%(refname:lstrip=2) sha:%(objectname) head:%(if)%(HEAD)%(then)Y%(else)N%(end) worktree:%(if)%(worktreepath)%(then)Y%(else)N%(end) symref:%(if)%(symref)%(then)Y%(else)N%(end) upstream:%(upstream:lstrip=2) track:%(upstream:track,nobracket) --sort=refname --include-root-refs HEAD refs/heads/ refs/remotes/ |
      |        | git config git-town-branch.feature.branchtype prototype                                                                                                                                                                                                                                                                                                                 |
      |        | git for-each-ref --format=refname:%(refname) branchname:%(refname:lstrip=2) sha:%(objectname) head:%(if)%(HEAD)%(then)Y%(else)N%(end) worktree:%(if)%(worktreepath)%(then)Y%(else)N%(end) symref:%(if)%(symref)%(then)Y%(else)N%(end) upstream:%(upstream:lstrip=2) track:%(upstream:track,nobracket) --sort=refname --include-root-refs HEAD refs/heads/ refs/remotes/ |
      |        | git config -lz --includes --global                                                                                                                                                                                                                                                                                                                                      |
      |        | git config -lz --includes --local                                                                                                                                                                                                                                                                                                                                       |
    And Git Town prints:
      """
      Ran 9 shell commands
      """
    And Git Town prints:
      """
      branch "feature" is now a prototype branch
      """
    And branch "feature" now has type "prototype"

  Scenario: undo
    When I run "git-town undo --verbose"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                                                                                                                                                                                                                                                                                                                                 |
      |        | git version                                                                                                                                                                                                                                                                                                                                                             |
      |        | git rev-parse --show-toplevel                                                                                                                                                                                                                                                                                                                                           |
      |        | git config -lz --includes --global                                                                                                                                                                                                                                                                                                                                      |
      |        | git config -lz --includes --local                                                                                                                                                                                                                                                                                                                                       |
      |        | git status -z --ignore-submodules                                                                                                                                                                                                                                                                                                                                       |
      |        | git rev-parse --verify -q MERGE_HEAD                                                                                                                                                                                                                                                                                                                                    |
      |        | git rev-parse --absolute-git-dir                                                                                                                                                                                                                                                                                                                                        |
      |        | git stash list                                                                                                                                                                                                                                                                                                                                                          |
      |        | git for-each-ref --format=refname:%(refname) branchname:%(refname:lstrip=2) sha:%(objectname) head:%(if)%(HEAD)%(then)Y%(else)N%(end) worktree:%(if)%(worktreepath)%(then)Y%(else)N%(end) symref:%(if)%(symref)%(then)Y%(else)N%(end) upstream:%(upstream:lstrip=2) track:%(upstream:track,nobracket) --sort=refname --include-root-refs HEAD refs/heads/ refs/remotes/ |
      |        | git remote get-url origin                                                                                                                                                                                                                                                                                                                                               |
      |        | git rev-parse --verify --abbrev-ref @{-1}                                                                                                                                                                                                                                                                                                                               |
      |        | git remote get-url origin                                                                                                                                                                                                                                                                                                                                               |
      |        | git config --unset git-town-branch.feature.branchtype                                                                                                                                                                                                                                                                                                                   |
    And Git Town prints:
      """
      Ran 13 shell commands
      """
    And branch "feature" now has type "feature"
