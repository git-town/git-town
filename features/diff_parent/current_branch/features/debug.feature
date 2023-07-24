Feature: display debug statistics

  Scenario: feature branch
    And the current branch is a feature branch "feature"
    When I run "git-town diff-parent --debug"
    Then it runs the commands
      | BRANCH  | TYPE     | COMMAND                         |
      |         | backend  | git config -lz --local          |
      |         | backend  | git config -lz --global         |
      |         | backend  | git rev-parse                   |
      |         | backend  | git rev-parse --show-toplevel   |
      |         | backend  | git version                     |
      |         | backend  | git branch -a                   |
      |         | backend  | git status                      |
      |         | backend  | git rev-parse --abbrev-ref HEAD |
      | feature | frontend | git diff main..feature          |
