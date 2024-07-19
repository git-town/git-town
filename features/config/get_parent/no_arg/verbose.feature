Feature: display the parent of a top-level feature branch

  Background:
    Given a Git repo clone
    And the branch
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    When I run "git-town config get-parent --verbose"

  Scenario: result
    Then it runs the commands
      | BRANCH | TYPE    | COMMAND                            |
      |        | backend | git version                        |
      |        | backend | git rev-parse --show-toplevel      |
      |        | backend | git config -lz --includes --global |
      |        | backend | git config -lz --includes --local  |
      |        | backend | git rev-parse --abbrev-ref HEAD    |
    And it prints:
      """
      Ran 5 shell commands.
      """
