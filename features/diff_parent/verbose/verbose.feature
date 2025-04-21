Feature: display all executed Git commands

  Scenario: feature branch
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    And the current branch is "feature"
    When I run "git-town diff-parent --verbose"
    Then Git Town runs the commands
      | BRANCH  | TYPE     | COMMAND                                          |
      |         | backend  | git version                                      |
      |         | backend  | git rev-parse --show-toplevel                    |
      |         | backend  | git config -lz --includes --global               |
      |         | backend  | git config -lz --includes --local                |
      |         | backend  | git status -z --ignore-submodules                |
      |         | backend  | git rev-parse -q --verify MERGE_HEAD             |
      |         | backend  | git rev-parse --absolute-git-dir                 |
      |         | backend  | git stash list                                   |
      |         | backend  | git -c core.abbrev=40 branch -vva --sort=refname |
      |         | backend  | git remote get-url origin                        |
      | feature | frontend | git diff main..feature                           |
    And Git Town prints:
      """
      Ran 11 shell commands.
      """
