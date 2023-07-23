Feature: display debug statistics

  Scenario: select another branch
    Given a feature branch "parent"
    And a feature branch "child" as a child of "parent"
    And the current branch is "child"
    When I run "git-town set-parent --debug" and answer the prompts:
      | PROMPT                                      | ANSWER        |
      | Please specify the parent branch of 'child' | [DOWN][ENTER] |
    Then it runs the debug commands
      | git config -lz --local                          |
      | git config -lz --global                         |
      | git rev-parse                                   |
      | git rev-parse --show-toplevel                   |
      | git version                                     |
      | git branch -a                                   |
      | git status                                      |
      | git rev-parse --abbrev-ref HEAD                 |
      | git config --unset git-town-branch.child.parent |
      | git branch                                      |
      | git config git-town-branch.child.parent main    |
    And this branch lineage exists now
      | BRANCH | PARENT |
      | child  | main   |
      | parent | main   |
