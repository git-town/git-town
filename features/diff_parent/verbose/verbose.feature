Feature: display all executed Git commands

  Scenario: feature branch
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    And the current branch is "feature"
    When I run "git-town diff-parent --verbose"
    Then Git Town runs the commands
      | BRANCH  | TYPE     | COMMAND                                                                                                                                                                                                                                                                                                                                          |
      |         | backend  | git version                                                                                                                                                                                                                                                                                                                                      |
      |         | backend  | git rev-parse --show-toplevel                                                                                                                                                                                                                                                                                                                    |
      |         | backend  | git config -lz --global                                                                                                                                                                                                                                                                                                                          |
      |         | backend  | git config -lz --local                                                                                                                                                                                                                                                                                                                           |
      |         | backend  | git config -lz                                                                                                                                                                                                                                                                                                                                   |
      |         | backend  | git status -z --ignore-submodules                                                                                                                                                                                                                                                                                                                |
      |         | backend  | git rev-parse --verify -q MERGE_HEAD                                                                                                                                                                                                                                                                                                             |
      |         | backend  | git rev-parse --absolute-git-dir                                                                                                                                                                                                                                                                                                                 |
      |         | backend  | git remote get-url origin                                                                                                                                                                                                                                                                                                                        |
      |         | backend  | git stash list                                                                                                                                                                                                                                                                                                                                   |
      |         | backend  | git for-each-ref "--format=refname:%(refname) branchname:%(refname:lstrip=2) sha:%(objectname) head:%(if)%(HEAD)%(then)Y%(else)N%(end) worktree:%(if)%(worktreepath)%(then)Y%(else)N%(end) symref:%(if)%(symref)%(then)Y%(else)N%(end) upstream:%(upstream:lstrip=2) track:%(upstream:track,nobracket)" --sort=refname refs/heads/ refs/remotes/ |
      |         | backend  | git remote                                                                                                                                                                                                                                                                                                                                       |
      | feature | frontend | git diff --merge-base main feature                                                                                                                                                                                                                                                                                                               |
    And Git Town prints:
      """
      Ran 13 shell commands.
      """
