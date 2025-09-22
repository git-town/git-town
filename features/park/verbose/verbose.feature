Feature: park the current branch verbosely

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    And the current branch is "feature"
    When I run "git-town park --verbose"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | TYPE    | COMMAND                                                                                                                                                                                                                                                                                                                                          |
      |        | backend | git version                                                                                                                                                                                                                                                                                                                                      |
      |        | backend | git rev-parse --show-toplevel                                                                                                                                                                                                                                                                                                                    |
      |        | backend | git config -lz --global                                                                                                                                                                                                                                                                                                                          |
      |        | backend | git config -lz --local                                                                                                                                                                                                                                                                                                                           |
      |        | backend | git config -lz                                                                                                                                                                                                                                                                                                                                   |
      |        | backend | git for-each-ref "--format=refname:%(refname) branchname:%(refname:lstrip=2) sha:%(objectname) head:%(if)%(HEAD)%(then)Y%(else)N%(end) worktree:%(if)%(worktreepath)%(then)Y%(else)N%(end) symref:%(if)%(symref)%(then)Y%(else)N%(end) upstream:%(upstream:lstrip=2) track:%(upstream:track,nobracket)" --sort=refname refs/heads/ refs/remotes/ |
      |        | backend | git config git-town-branch.feature.branchtype parked                                                                                                                                                                                                                                                                                             |
      |        | backend | git for-each-ref "--format=refname:%(refname) branchname:%(refname:lstrip=2) sha:%(objectname) head:%(if)%(HEAD)%(then)Y%(else)N%(end) worktree:%(if)%(worktreepath)%(then)Y%(else)N%(end) symref:%(if)%(symref)%(then)Y%(else)N%(end) upstream:%(upstream:lstrip=2) track:%(upstream:track,nobracket)" --sort=refname refs/heads/ refs/remotes/ |
      |        | backend | git config -lz --global                                                                                                                                                                                                                                                                                                                          |
      |        | backend | git config -lz --local                                                                                                                                                                                                                                                                                                                           |
    And Git Town prints:
      """
      Ran 10 shell commands
      """
    And Git Town prints:
      """
      branch "feature" is now parked
      """
    And branch "feature" now has type "parked"

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
      |        | backend | git config git-town-branch.feature.branchtype feature                                                                                                                                                                                                                                                                                            |
    And Git Town prints:
      """
      Ran 14 shell commands
      """
    And branch "feature" now has type "feature"
