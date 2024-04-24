Feature: display all executed Git commands

  @debug @this
  Scenario: select another branch
    Given a feature branch "parent"
    And a feature branch "child" as a child of "parent"
    And the current branch is "child"
    When I run "git-town set-parent --verbose" and enter into the dialog:
      | DIALOG                 | KEYS     |
      | parent branch of child | up enter |
    Then it runs the commands
      | BRANCH | TYPE    | COMMAND                                         |
      |        | backend | git version                                     |
      |        | backend | git config -lz --global                         |
      |        | backend | git config -lz --local                          |
      |        | backend | git rev-parse --show-toplevel                   |
      |        | backend | git status --long --ignore-submodules           |
      |        | backend | git stash list                                  |
      |        | backend | git branch -vva --sort=refname                  |
      |        | backend | git config --unset git-town-branch.child.parent |
      |        | backend | git config -lz --global                         |
      |        | backend | git config -lz --local                          |
      |        | backend | git config git-town-branch.child.parent main    |
      |        | backend | git config -lz --global                         |
      |        | backend | git config -lz --local                          |
    And it prints:
      """
      Ran 13 shell commands.
      """
    And this lineage exists now
      | BRANCH | PARENT |
      | child  | main   |
      | parent | main   |
