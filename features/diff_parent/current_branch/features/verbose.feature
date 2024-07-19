Feature: display all executed Git commands

  Scenario: feature branch
    Given a Git repo clone
    And the branch
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    And the current branch is "feature"
    When I run "git-town diff-parent --verbose"
    Then it runs the commands
      | BRANCH  | TYPE     | COMMAND                               |
      |         | backend  | git version                           |
      |         | backend  | git rev-parse --show-toplevel         |
      |         | backend  | git config -lz --includes --global    |
      |         | backend  | git config -lz --includes --local     |
      |         | backend  | git status --long --ignore-submodules |
      |         | backend  | git stash list                        |
      |         | backend  | git branch -vva --sort=refname        |
      | feature | frontend | git diff main..feature                |
    And it prints:
      """
      Ran 8 shell commands.
      """
