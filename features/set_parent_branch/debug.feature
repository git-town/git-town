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
      |        | backend | git version                                     |
      |        | backend | git config -lz --local                          |
      |        | backend | git config -lz --global                         |
      |        | backend | git rev-parse --show-toplevel                   |
      |        | backend | git rev-parse --show-toplevel                   |
      |        | backend | git branch -vva                                 |
      |        | backend | git branch -a                                   |
      |        | backend | git config --unset git-town-branch.child.parent |
      |        | backend | git branch                                      |
      |        | backend | git config git-town-branch.child.parent main    |
    And it prints:
      """
      Ran 10 shell commands.
      """
    And this branch lineage exists now
      | BRANCH | PARENT |
      | child  | main   |
      | parent | main   |
