Feature: display the parent of a top-level feature branch

  Background:
    Given the current branch is a feature branch "feature"
    When I run "git-town config get-parent --verbose"

  @this
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
