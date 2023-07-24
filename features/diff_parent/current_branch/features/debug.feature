Feature: display debug statistics

  Scenario: feature branch
    And the current branch is a feature branch "feature"
    When I run "git-town diff-parent --debug"
    Then it runs the commands
      | BRANCH  | TYPE     | COMMAND                       |
      |         | backend  | git version                   |
      |         | backend  | git config -lz --local        |
      |         | backend  | git config -lz --global       |
      |         | backend  | git rev-parse --show-toplevel |
      |         | backend  | git branch -vva               |
      |         | backend  | git branch -a                 |
      |         | backend  | git rev-parse --show-toplevel |
      | feature | frontend | git diff main..feature        |
    And it prints:
      """
      Ran 8 shell commands.
      """
