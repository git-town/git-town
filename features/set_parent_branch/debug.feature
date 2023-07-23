Feature: display debug statistics

  Scenario: select another branch
    Given a feature branch "parent"
    And a feature branch "child" as a child of "parent"
    And the current branch is "child"
    When I run "git-town set-parent --debug" and answer the prompts:
      | PROMPT                                      | ANSWER        |
      | Please specify the parent branch of 'child' | [DOWN][ENTER] |
    Then it runs the commands
      | BRANCH | TYPE    | COMMAND                                         |
      |        | backend | git config -lz --local                          |
      |        | backend | git config -lz --global                         |
      |        | backend | git rev-parse                                   |
      |        | backend | git rev-parse --show-toplevel                   |
      |        | backend | git version                                     |
      |        | backend | git branch -a                                   |
      |        | backend | git status                                      |
      |        | backend | git rev-parse --abbrev-ref HEAD                 |
      |        | backend | git config --unset git-town-branch.child.parent |
      |        | backend | git branch                                      |
      |        | backend | git config git-town-branch.child.parent main    |
    And this branch lineage exists now
      | BRANCH | PARENT |
      | child  | main   |
      | parent | main   |
