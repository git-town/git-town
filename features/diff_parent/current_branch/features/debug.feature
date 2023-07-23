Feature: display debug statistics

  Scenario: feature branch
    And the current branch is a feature branch "feature"
    When I run "git-town diff-parent --debug"
    Then it runs the debug commands
      | git config -lz --local          |
      | git config -lz --global         |
      | git rev-parse                   |
      | git rev-parse --show-toplevel   |
      | git version                     |
      | git branch -a                   |
      | git status                      |
      | git rev-parse --abbrev-ref HEAD |
